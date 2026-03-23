package core

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func LoginQR(client *whatsmeow.Client, log waLog.Logger) error {
	if client.Store.ID != nil {
		return client.Connect()
	}

	qrChan, err := client.GetQRChannel(context.Background())
	if err != nil {
		return fmt.Errorf("gagal get QR channel: %w", err)
	}

	if err = client.Connect(); err != nil {
		return fmt.Errorf("gagal connect: %w", err)
	}

	for evt := range qrChan {
		switch evt.Event {
		case "code":
			fmt.Fprintln(os.Stdout, "")
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			fmt.Fprintln(os.Stdout, "")
		case whatsmeow.QRChannelSuccess.Event:
			log.Infof("Login berhasil!")
		case whatsmeow.QRChannelTimeout.Event:
			return fmt.Errorf("QR timeout, jalankan ulang bot")
		case whatsmeow.QRChannelErrUnexpectedEvent.Event:
			return fmt.Errorf("unexpected event saat QR login")
		case whatsmeow.QRChannelClientOutdated.Event:
			return fmt.Errorf("versi client outdated")
		case "error":
			return fmt.Errorf("QR error: %v", evt.Error)
		}
	}
	return nil
}

func LoginPairCode(client *whatsmeow.Client, phone string, log waLog.Logger) error {
	if client.Store.ID != nil {
		return client.Connect()
	}

	if err := client.Connect(); err != nil {
		return fmt.Errorf("gagal connect: %w", err)
	}

	time.Sleep(2 * time.Second)

	code, err := client.PairPhone(
		context.Background(),
		phone,
		true,
		whatsmeow.PairClientChrome,
		"Chrome (Linux)",
	)
	if err != nil {
		return fmt.Errorf("gagal pair phone: %w", err)
	}

	fmt.Printf("\n[ KODE PAIRING: %s ]\n\n", code)
	fmt.Println("WhatsApp > Pengaturan > Perangkat Tertaut > Tautkan dengan Nomor Telepon")
	fmt.Println("Masukkan kode di atas, tunggu hingga terhubung...")

	done := make(chan struct{})
	id := client.AddEventHandler(func(evt interface{}) {
		if _, ok := evt.(*events.PairSuccess); ok {
			close(done)
		}
	})
	defer client.RemoveEventHandler(id)
	<-done
	return nil
}
