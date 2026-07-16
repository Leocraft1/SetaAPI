#!/usr/bin/env python3
"""
Script di migrazione: JSON (linee + percorsi) -> MySQL/MariaDB.

=== STRUTTURA DEI FILE JSON (verificata sui file forniti) ===

1) File principale delle linee (es. "rc_new.json"):
   [
     {"linea": "1", "codes": ["MO1-As-150", "MO1-As-139", ...]},
     {"linea": "9", "codes": [...]},
     ...
   ]
   NB: "linea" NON e' sempre un numero puro (es. "14 (2025)", "rossa",
   "verde", "C", "5taxi"), quindi viene trattato come stringa/VARCHAR.

2) File di percorso, uno per ogni codice in "codes"
   (es. "MO1-As-139.json", nella cartella --percorsi-dir):
   [
     {"desc": "POLO LEONARDO", "code": "MO6783", "islast": false},
     {"desc": "POLO LEONARDO 1", "code": "MO2928", "islast": false},
     ...
     {"desc": "MARINUZZI", "code": "MO2703", "islast": true}
   ]
   -> lista ORDINATA di fermate. "code" e' il codice fermata (es. "MO6783")
      e corrisponde alla colonna `fermate.codice` (PK della tabella
      `fermate` gia' esistente, gestita dalla tua API - NON viene toccata
      da questo script). "islast" indica il capolinea del percorso.

=== SCHEMA DB ATTESO ===

    CREATE TABLE linee (
        id   VARCHAR(50) PRIMARY KEY,   -- es. "1", "14 (2025)", "rossa"
        nome VARCHAR(255)
    );

    CREATE TABLE percorsi (
        id       VARCHAR(50) PRIMARY KEY,  -- es. "MO1-As-139"
        linea_id VARCHAR(50) NOT NULL,
        FOREIGN KEY (linea_id) REFERENCES linee(id)
    );

    CREATE TABLE percorso_fermate (
        percorso_id VARCHAR(50) NOT NULL,
        fermata_id  VARCHAR(50) NOT NULL,   -- referenzia fermate.codice
        ordine      INT NOT NULL,
        capolinea   TINYINT(1) NOT NULL DEFAULT 0,
        PRIMARY KEY (percorso_id, ordine),
        FOREIGN KEY (percorso_id) REFERENCES percorsi(id),
        FOREIGN KEY (fermata_id) REFERENCES fermate(codice)
    );

    CREATE INDEX idx_percorso_fermate_percorso ON percorso_fermate (percorso_id, ordine);

=== USO ===
    pip install pymysql
    python migrate_json_to_db.py \
        --host localhost --user root --password *** --database trasporti \
        --linee-file ./rc_new.json --percorsi-dir ./percorsi \
        --dry-run

Rilancia senza --dry-run quando i risultati ti convincono. Lo script e'
idempotente: puoi rilanciarlo piu' volte senza creare duplicati.
"""

import argparse
import json
import logging
import re
import sys
from pathlib import Path

import pymysql
from pymysql.cursors import DictCursor

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
)
log = logging.getLogger("migrate")

# Passo tra un ordine e il successivo: lascia "buchi" per poter inserire
# in futuro una fermata in mezzo senza rinumerare tutto il percorso.
ORDINE_STEP = 10


# ---------------------------------------------------------------------------
# LETTURA DATI DAI JSON
# ---------------------------------------------------------------------------

def carica_linee(linee_file: Path) -> list[dict]:
    """Ritorna la lista [{'linea': str, 'codes': [str, ...]}, ...]"""
    with open(linee_file, encoding="utf-8") as f:
        return json.load(f)


# Alcuni codici in rc_new.json hanno un'annotazione finale tra parentesi
# (es. "MO3-As-343 (2025)", "MO10-Di-1041 (12-2025)") che pero' NON fa
# parte del nome del file fisico su disco (es. "MO3-As-343.json").
# L'annotazione viene mantenuta nell'id del percorso salvato nel DB
# (per non perdere l'informazione e per evitare collisioni tra codici
# come "MO10-Di-1041" e "MO10-Di-1041 (12-2025)"), ma viene rimossa solo
# per individuare il file da leggere.
_RE_ANNOTAZIONE = re.compile(r"\s*\([^)]*\)\s*$")


def nome_file_percorso(codice_percorso: str) -> str:
    return _RE_ANNOTAZIONE.sub("", codice_percorso).strip()


def carica_percorso(percorsi_dir: Path, codice_percorso: str) -> list[dict]:
    """Ritorna la lista ordinata di fermate di un percorso:
    [{'desc': str, 'code': str, 'islast': bool}, ...]
    """
    nome_file = nome_file_percorso(codice_percorso)
    path = percorsi_dir / f"{nome_file}.json"
    if not path.exists():
        raise FileNotFoundError(
            f"File percorso non trovato: {path} "
            f"(derivato dal codice '{codice_percorso}')"
        )
    with open(path, encoding="utf-8") as f:
        return json.load(f)


# ---------------------------------------------------------------------------
# LOGICA DI MIGRAZIONE
# ---------------------------------------------------------------------------

def raccogli_codici_fermata(linee: list[dict], percorsi_dir: Path) -> set[str]:
    """Prima passata: legge tutti i file di percorso e raccoglie i codici
    fermata usati, per poter verificare in un'unica query quali esistono
    gia' nella tabella `fermate`."""
    codici = set()
    for linea in linee:
        for codice_percorso in linea["codes"]:
            fermate = carica_percorso(percorsi_dir, codice_percorso)
            for f in fermate:
                codici.add(f["code"])
    return codici


def fermate_esistenti(conn, codici: set[str]) -> set[str]:
    """Ritorna il sottoinsieme di `codici` presente in fermate.codice."""
    if not codici:
        return set()
    codici = list(codici)
    trovate = set()
    with conn.cursor() as cur:
        # IN (...) a blocchi per evitare limiti su query troppo lunghe
        CHUNK = 1000
        for i in range(0, len(codici), CHUNK):
            blocco = codici[i:i + CHUNK]
            placeholders = ",".join(["%s"] * len(blocco))
            cur.execute(
                f"SELECT codice FROM fermate WHERE codice IN ({placeholders})",
                blocco,
            )
            trovate.update(row[0] for row in cur.fetchall())
    return trovate


def migra(conn, linee_file: Path, percorsi_dir: Path, dry_run: bool = False):
    linee = carica_linee(linee_file)
    log.info("Trovate %d linee nel file %s", len(linee), linee_file)

    log.info("Prima passata: raccolgo i codici fermata usati nei percorsi...")
    codici_usati = raccogli_codici_fermata(linee, percorsi_dir)
    codici_ok = fermate_esistenti(conn, codici_usati)
    codici_mancanti = codici_usati - codici_ok
    if codici_mancanti:
        log.warning(
            "%d codici fermata usati nei JSON NON esistono in `fermate` "
            "e verranno saltati: %s",
            len(codici_mancanti),
            ", ".join(sorted(codici_mancanti)[:20])
            + (" ..." if len(codici_mancanti) > 20 else ""),
        )

    tot_percorsi = 0
    tot_righe_inserite = 0
    tot_righe_saltate = 0

    with conn.cursor(DictCursor) as cur:
        for linea in linee:
            linea_id = linea["linea"]

            if not dry_run:
                cur.execute(
                    """
                    INSERT INTO linee (id, nome)
                    VALUES (%s, %s)
                    ON DUPLICATE KEY UPDATE nome = VALUES(nome)
                    """,
                    (linea_id, linea_id),
                )

            for codice_percorso in linea["codes"]:
                tot_percorsi += 1
                fermate = carica_percorso(percorsi_dir, codice_percorso)

                if dry_run:
                    valide = [f for f in fermate if f["code"] in codici_ok]
                    tot_righe_inserite += len(valide)
                    tot_righe_saltate += len(fermate) - len(valide)
                    continue

                cur.execute(
                    """
                    INSERT INTO percorsi (id, linea_id)
                    VALUES (%s, %s)
                    ON DUPLICATE KEY UPDATE linea_id = VALUES(linea_id)
                    """,
                    (codice_percorso, linea_id),
                )

                # Ripulisce le associazioni esistenti per questo percorso,
                # cosi' un rilancio riflette fedelmente il contenuto del JSON
                # anche se ordine o fermate sono cambiate nel frattempo.
                cur.execute(
                    "DELETE FROM percorso_fermate WHERE percorso_id = %s",
                    (codice_percorso,),
                )

                righe = []
                ordine = 0
                for f in fermate:
                    if f["code"] not in codici_ok:
                        tot_righe_saltate += 1
                        continue
                    righe.append((
                        codice_percorso,
                        f["code"],
                        ordine,
                        1 if f.get("islast") else 0,
                    ))
                    ordine += ORDINE_STEP

                if righe:
                    cur.executemany(
                        """
                        INSERT INTO percorso_fermate
                            (percorso_id, fermata_id, ordine, capolinea)
                        VALUES (%s, %s, %s, %s)
                        """,
                        righe,
                    )
                tot_righe_inserite += len(righe)

    if dry_run:
        log.info(
            "[DRY RUN] Nessuna scrittura eseguita. Percorsi: %d, "
            "associazioni valide: %d, saltate per fermata mancante: %d",
            tot_percorsi, tot_righe_inserite, tot_righe_saltate,
        )
    else:
        conn.commit()
        log.info(
            "Migrazione completata. Percorsi: %d, associazioni inserite: %d, "
            "saltate per fermata mancante: %d",
            tot_percorsi, tot_righe_inserite, tot_righe_saltate,
        )


# ---------------------------------------------------------------------------
# ENTRY POINT
# ---------------------------------------------------------------------------

def main():
    parser = argparse.ArgumentParser(
        description=__doc__,
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    parser.add_argument("--host", required=True)
    parser.add_argument("--port", type=int, default=3306)
    parser.add_argument("--user", required=True)
    parser.add_argument("--password", required=True)
    parser.add_argument("--database", required=True)
    parser.add_argument("--linee-file", required=True, type=Path,
                         help="Path al file JSON tipo rc_new.json")
    parser.add_argument("--percorsi-dir", required=True, type=Path,
                         help="Cartella contenente i file JSON di ogni "
                              "percorso (es. MO1-As-139.json)")
    parser.add_argument("--dry-run", action="store_true",
                         help="Simula la migrazione senza scrivere nulla nel DB")
    args = parser.parse_args()

    conn = pymysql.connect(
        host=args.host,
        port=args.port,
        user=args.user,
        password=args.password,
        database=args.database,
        charset="utf8mb4",
        autocommit=False,
    )

    try:
        migra(conn, args.linee_file, args.percorsi_dir, dry_run=args.dry_run)
    except Exception:
        conn.rollback()
        log.exception("Migrazione fallita, rollback eseguito")
        sys.exit(1)
    finally:
        conn.close()


if __name__ == "__main__":
    main()
