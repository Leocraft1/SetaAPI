package gtfs

import (
	"fmt"
	"log"
	"setaapi/config"
	"setaapi/internal/repository"
	"setaapi/internal/repository/static"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	gtfsparserwr "github.com/Leocraft1/gtfsparser-with-reader"
)

func UpdateStatic() error {
	return repository.WithAdvisoryLock(repository.DB_STATIC, "setaapi:update-static-gtfs", updateStatic)
}

func updateStatic() error {
	// start := time.Now() // Uncomment to measure import duration.
	log.Println("Started static GTFS update.")
	feed, err := getStatic()
	if err != nil {
		log.Fatal("UpdateStatic error: ", err)
	}
	if err := static.SaveSnapshot(feed); err != nil {
		return err
	}
	// log.Printf("Static GTFS updated in %s", time.Since(start))
	return nil
}

func getStatic() (*gtfsparserwr.Feed, error) {
	// 1. Configure SSH Config
    sshConfig := &ssh.ClientConfig{
        User: config.FTP_USERNAME,
        Auth: []ssh.AuthMethod{
            ssh.Password(config.FTP_PASSWORD), // Or use PublicKeys(...)
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Use strict checking in production
    }

    // 2. Connect to the SSH server
    conn, err := ssh.Dial("tcp", config.SETA_GTFS_FTP, sshConfig)
    if err != nil {
        panic(fmt.Sprintf("Failed to dial: %s", err))
    }
    defer conn.Close()

    // 3. Create the SFTP client
    client, err := sftp.NewClient(conn)
    if err != nil {
        panic(fmt.Sprintf("Failed to create SFTP client: %s", err))
    }
    defer client.Close()

	zipped, err := client.Open("seta.zip")
	if err != nil {
		return nil, fmt.Errorf("Unable to open static GTFS file")
	}

	feed := gtfsparserwr.NewFeed()

	if err := feed.ParseReader(zipped); err != nil {
		return nil, fmt.Errorf("GetStaticFeed: parse gtfs: %w", err)
	}

	return feed, nil 
}
