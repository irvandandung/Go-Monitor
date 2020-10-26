package collector

import (
        "bufio"
        "log"
        "os"
        "strings"
)

// struct untuk menyiman nilai idle sebelumnya dan total
// sebelumnya
type CPUStats struct {
        PrevIdleTime  uint64
        PrevTotalTime uint64
}
 
func GetCPU(stat *CPUStats) float64 {
        file, err := os.Open("/proc/stat")
        if err != nil {
                log.Fatal(err)
        }
        defer file.Close()
        // parsing isi /proc/stat
        return ParseCPU(file, stat)
}
 
func ParseCPU(file *os.File, stat *CPUStats) float64 {
        var cpuStats float64
 
        scanner := bufio.NewScanner(file)
        // hanya ambil baris pertama
        scanner.Scan()
        // cth: cpu  1679009 37112 377610 11558573 16807 0 19461 0 0 0
        // potong 5 karakter "cpu  "
        parts := strings.Fields(scanner.Text()[5:])
        // set bagian ke empat "11558573" ke variable idleTime
        // set totalTime 0, dan nanti akan terisi hasil penjumlahan
        // semua elemen dari parts
        idleTime := toUint64(parts[3])
        totalTime := uint64(0)
 
        for _, x := range parts {
                u := toUint64(x)
                // jumlahkan semua bagian pada parts
                totalTime = totalTime + u
        }
 
        // hitung selisih antara nilai idle sekarang dan sebelumnya
        // dan hitung selisih total sekarang dan total sebelumnya
        deltaIdleTime := idleTime - stat.PrevIdleTime
        deltaTotalTime := totalTime - stat.PrevTotalTime
 
        // perbarui nilai idle sebelumnya menjadi nilai idle sekarang
        // dan perbarui nilai total sebelumnya menjadi nilai total sekarang
        stat.PrevIdleTime = idleTime
        stat.PrevTotalTime = totalTime
 
        // hitung persentase
        cpuStats = (1.0 - float64(deltaIdleTime)/float64(deltaTotalTime)) * 100.0
        // kembalikan nilai persentase cpu
        return cpuStats
}