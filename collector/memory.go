package collector

import (
    "bufio"
    "os"
    "regexp"
    "strings"
)

var MemStats = regexp.MustCompile("^(MemTotal|MemFree|MemAvailable|Buffers|Cached)$")

func GetMemory() (map[string]uint64, error) {
    // buka berkas /proc/meminfo
    file, err := os.Open("/proc/meminfo")
    if err != nil {
        return nil, err
    }
    defer file.Close()
    // kemalikan hasil parsing dari /proc/meminfo
    return ParseMemory(file)
}

func ParseMemory(file *os.File) (map[string]uint64, error) {
    // buat map baru map[string]uint64 yang akan berisi informasi
    // mengenai memori
    meminfo := make(map[string]uint64)
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        // pisah setiap bagian dengan spasi sebagai pemisah
        parts := strings.Fields(scanner.Text())
        // hilangkan ":" pada bagian pertama, dan set ke variable key
        key := strings.TrimRight(parts[0], ":")
        if MemStats.Match([]byte(key)) {
            // jika key sesuai dengan string yang berada pada variable
            // Memstats, tambahkan elemen map baru dengan key key dengan
            // value kolom kedua * 1024
            meminfo[key] = toUint64(parts[1]) * 1024
        }
    }
 
    // karena penggunaan memori tidak tercatat di /proc/meminfo, kita
    // lakukan penghitungan secara manual
    meminfo["Used"] = meminfo["MemTotal"] - meminfo["MemFree"] - meminfo["Buffers"] - meminfo["Cached"]
 
    // kembalikan map yang berisi informasi mengenai memori
    return meminfo, nil
}