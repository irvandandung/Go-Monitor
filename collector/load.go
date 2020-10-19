package collector
 
import (
    "io/ioutil"
    "strings"
)
 
func GetLoad() (map[string]float64, error) {
    // baca berkas /proc/loadavg
    file, err := ioutil.ReadFile("/proc/loadavg")
    if err != nil {
        return nil, err
    }
    // parsing isi berkas /proc/loadavg dengan parseLoad(), setelah
    // diparsing akan mengembalikan load dalam bentuk map[string]float64
    // dan untuk error kembalikan nil
    return ParseLoad(string(file)), nil
}
 
func ParseLoad(data string) map[string]float64 {
    // inisialisasi map baru dengan tipe map[string]float64
    loads := make(map[string]float64)
    // pisahkan setiap bagian string yang didapat dari getLoad()
    // dengan spasi sebagai pemisah, setelah dipisah akan menjadi
    // slice yang berisi substring
    parts := strings.Fields(data)
    // masukan data dari slice parts ke map, karena hanya 3 data yang
    // akan kita ambil, jadi kita hanya menggunakan parts slice dengan
    // index 0, 1, dan 2
    loads["load1"] = toFloat(parts[0])
    loads["load5"] = toFloat(parts[1])
    loads["load15"] = toFloat(parts[2])
    return loads
}