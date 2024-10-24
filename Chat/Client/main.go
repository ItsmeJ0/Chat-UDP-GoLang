package main

import (
	"bufio"   // Untuk membaca input dari pengguna
	"fmt"     // Untuk mencetak ke layar
	"net"     // Untuk menghubungkan melalui jaringan
	"os"      // Untuk membaca input dari command line
	"strings" // Untuk manipulasi string, seperti trimming newline
)

func main() {
	// Menentukan alamat server (localhost di port 8080)
	serverAddr := net.UDPAddr{
		Port: 8080,                     // Port 8080 yang akan digunakan
		IP:   net.ParseIP("127.0.0.1"), // Alamat IP localhost
	}

	// Menghubungkan ke server
	conn, err := net.DialUDP("udp", nil, &serverAddr) // Membuat koneksi UDP ke server
	if err != nil {                                   // Jika ada error saat mencoba koneksi
		fmt.Println("Error:", err) // Tampilkan error dan keluar
		return
	}
	defer conn.Close() // Menutup koneksi setelah fungsi `main` selesai

	// Membaca input nama pengguna dari stdin (command line)
	reader := bufio.NewReader(os.Stdin)      // Membuat objek reader untuk membaca input dari command line
	fmt.Print("Enter your name: ")           // Menampilkan prompt untuk meminta nama pengguna
	nama_user, _ := reader.ReadString('\n')  // Membaca nama yang dimasukkan oleh pengguna
	nama_user = strings.TrimSpace(nama_user) // Menghapus karakter newline di akhir input

	// Mengirim nama ke server sebagai identitas pengguna
	_, err = conn.Write([]byte("NAME:" + nama_user)) // Mengirim nama dengan format "NAME:<nama>"
	if err != nil {                                  // Jika terjadi kesalahan saat mengirim nama
		fmt.Println("Error sending name:", err) // Tampilkan error
		return
	}

	// Goroutine untuk menerima pesan dari server secara terus menerus
	go func() {
		for { // Loop untuk terus-menerus mendengarkan pesan dari server
			buffer := make([]byte, 1024)          // Membuat buffer untuk menampung pesan
			n, _, err := conn.ReadFromUDP(buffer) // Membaca pesan dari server
			if err != nil {                       // Jika ada error saat menerima pesan
				fmt.Println("Error receiving message:", err)
				continue // Lewati iterasi jika error
			}

			// Menampilkan pesan yang diterima dari server
			fmt.Printf("\n%s\n", string(buffer[:n])) // Menampilkan pesan
		}
	}()

	// Loop untuk mengirim pesan
	for {
		// Membaca input pesan dari pengguna
		pesan_user, _ := reader.ReadString('\n')   // Membaca input pesan dari command line
		pesan_user = strings.TrimSpace(pesan_user) // Menghapus karakter newline di akhir pesan

		// Jika pengguna mengetik 'exit', maka keluar dari loop dan program akan selesai
		if pesan_user == "exit" {
			fmt.Println("Closing connection...") // Menampilkan pesan bahwa koneksi akan ditutup
			break                                // Keluar dari loop
		}

		fmt.Println("[you]: " + pesan_user)
		// Mengirim pesan ke server dengan format: "PESAN:<pesan>"
		_, err = conn.Write([]byte("PESAN:" + pesan_user)) // Mengirim pesan ke server
		if err != nil {                                    // Jika ada error saat mengirim pesan
			fmt.Println("Error sending message:", err)
			continue // Lanjutkan ke iterasi berikutnya jika ada error
		}
	}
}
