package collector

import (
    "bufio"
    "encoding/json"
    "os"
    "regexp"
    "strings"
    "syscall"
)

type StorageInfo struct {
    Size  uint64 `json:"size"`
    Avail uint64 `json:"avail"`
    Used  uint64 `json:"used"`
}

// variable yang berisi mountpoints yang akan diabaikan
var IgnoredFSTypes = regexp.MustCompile("^(autofs|binfmt_misc|bpf|cgroup2?|configfs|debugfs|devpts|devtmpfs|fusectl|hugetlbfs|iso9660|mqueue|nsfs|overlay|proc|procfs|pstore|rpc_pipefs|securityfs|selinuxfs|squashfs|sysfs|tracefs)$")

func GetStorage() (map[string]string, error){
	// buat map baru untuk menyimpan mountpoint sebagai key, dan size
    // dan avail sebagai value
    mp := make(map[string]string)
    // panggil GetMountPoints() untuk mendapatkan list mountpoint
    mountpoints, err := GetMountPoints()
    if err != nil {
        return nil, err
    }
    //
    for _, mountpoint := range mountpoints {
        // buat alias untuk struct syscall.Statsfs_t
        var stat syscall.Statfs_t
        // gunakan syscall statfs() untuk mendapatkan informasi
        // mengenakan mountpoint yang kita sertakan, informasi akan
        // disimpan pada struct stat (syscall.Statfs_t)
        syscall.Statfs(mountpoint, &stat)
        // kita akan mendapatkan nilai dari size, avail dan used dengan
        // dengan mengkalikan ukuran block dengan optimal transfer block
        // size stat.Bsize
        size := stat.Blocks * uint64(stat.Bsize)
        avail := stat.Bavail * uint64(stat.Bsize)
        // untuk mendapatkan kapasitas peyimpanan yang digunakan, kurangi
        // kapasitas penyimpan dengan kapasitas tersedia (size - avail)
        used := size - avail
        // simpan informasi size, avail, dan used pada struct StorageInfo
        mpstat := StorageInfo{
            Size:  size,
            Avail: avail,
            Used:  used,
        }
        // kita akan marshal struct StorageInfo/mpstat menjadi json
        // cth: {"size":405044625408,"free":202509873152,"avail":181863309312}
        m, _ := json.Marshal(mpstat)
        // simpan mountpoint sebagai key kedalam map mp dan m sebagai value.
        // karena m berbentuk []byte kita konversi dulu kedalam bentuk string
        mp[mountpoint] = string(m)
    }
    // kembalikan map mp
    return mp, nil
}

func GetMountPoints() ([]string, error){
	// buka berkas /proc/mounts untuk mengambil semua mountpoint
    file, err := os.Open("/proc/mounts")
    if err != nil {
        return nil, err
    }
    defer file.Close()
 
    // kembalikan hasil parsing dari /proc/mounts
    return ParseMountPoints(file)
}

func ParseMountPoints(file *os.File) ([]string, error){
	// buat slice yang nantinya akan diisi list mountpoint
    var mountpoints []string
    // pindai setiap baris pada berkas /proc/mounts
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        // pisah setiap bagian dengan spasi sebagai pemisah
        parts := strings.Fields(scanner.Text())
        // cek apakah kolom ketiga mengandung string di variable
        // IgnoredFSTypes
        if !IgnoredFSTypes.Match([]byte(parts[2])) {
            // jika tidak, kita tambahkan kolom kedua yang berisi
            // mountpoint ke slice mountpoints,
            mountpoints = append(mountpoints, parts[1])
        }
    }
    // kembalikan mountpoints yang telah terisi list mountpoint
    return mountpoints, nil
}