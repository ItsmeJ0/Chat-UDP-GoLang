package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	// Menentukan alamat server (localhost di port 8080)
	addr := net.UDPAddr{
		Port: 8080,                     // Port yang digunakan server adalah 8080
		IP:   net.ParseIP("127.0.0.1"), // Server mendengarkan pada alamat IP localhost (127.0.0.1)
	}

	// Membuat koneksi untuk mendengarkan di alamat tersebut
	conn, err := net.ListenUDP("udp", &addr) // Membuka koneksi UDP untuk mendengarkan di port 8080
	if err != nil {
		fmt.Println("Error:", err) // Jika terjadi kesalahan saat membuka koneksi, tampilkan error
		return
	}
	defer conn.Close() // Menutup koneksi saat program berakhir

	fmt.Println("Server listening on port 8080")

	// Map untuk menyimpan alamat client dan nama mereka
	clients := make(map[string]string) // Menyimpan pasangan alamat client dan nama
	buffer := make([]byte, 1024)       // Buffer untuk menyimpan pesan yang diterima (maksimum 1024 byte)

	for {
		// Membaca data dari client
		n, clientAddr, err := conn.ReadFromUDP(buffer) // Menerima data dari client, menyimpan alamat client ke clientAddr
		if err != nil {
			fmt.Println("Error receiving from client:", err) // Menampilkan error jika gagal menerima data
			continue
		}

		// Mengonversi data dari buffer menjadi string
		message := string(buffer[:n])    // Mengubah data yang diterima (byte) menjadi string
		clientKey := clientAddr.String() // Mengubah alamat client ke dalam string

		// Cek apakah pesan adalah pengiriman nama pengguna
		if strings.HasPrefix(message, "NAME:") { // Jika pesan diawali dengan "NAME:", itu berarti client mengirimkan nama pengguna
			name := strings.TrimPrefix(message, "NAME:")                // Menghapus awalan "NAME:" untuk mendapatkan nama pengguna
			clients[clientKey] = name                                   // Menyimpan nama pengguna yang dikirimkan sesuai dengan alamat client
			fmt.Printf("Client %s registered as %s\n", clientKey, name) // Menampilkan bahwa client telah terdaftar
			continue
		}

		// Jika pesan adalah pesan biasa
		if strings.HasPrefix(message, "PESAN:") { // Jika pesan diawali dengan "PESAN:", berarti client mengirimkan pesan
			pesan := strings.TrimPrefix(message, "PESAN:") // Menghapus awalan "PESAN:" untuk mendapatkan isi pesan
			name := clients[clientKey]                     // Mendapatkan nama client berdasarkan alamatnya

			// Broadcast pesan ke semua client, termasuk pengirim
			for addr := range clients { // Melakukan iterasi ke setiap alamat client yang tersimpan
				if addr != clientKey { // Handle agar tidak mengirim ke pengirim
					// Membuat pesan dengan format yang diinginkan
					broadcastMsg := fmt.Sprintf("[%s]: %s", name, pesan)     // Menggunakan format "[nama_pengirim]: pesan"
					udpAddr, _ := net.ResolveUDPAddr("udp", addr)            // Mengonversi alamat client ke format UDP
					_, err := conn.WriteToUDP([]byte(broadcastMsg), udpAddr) // Mengirim pesan ke client
					if err != nil {
						fmt.Println("Error sending message to client:", err) // Jika ada error saat mengirim pesan, tampilkan error
					}
				}
			}

			// Menampilkan pesan dari pengirim di server
			fmt.Printf("Message from %s: %s\n", name, pesan) // Menampilkan pesan dari client di server
		}

	}
}
