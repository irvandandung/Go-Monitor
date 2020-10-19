package main
 
import (
        "fmt"
        "log"
        "time"
        // import package collector sesuai dengan module path
        "github.com/irvandandung/gomon/collector"
        // import package command sesuai dengan module path
        "github.com/irvandandung/gomon/pkg/models/command"
)
 
func main() {
        // inisialisasi database dengan memanggil influxConnect()
        // db := InfluxConnect()
        // defer db.Close()

        
        // kita akan melakukan looping secara terus-menerus
        for {
                // panggil fungsi getLoad() yang saat ini berada dalam
                // package collector untuk mendapatkan informasi load yang
                // telah diparsing, dan lakukan pengecekan, jika terdapat error
                // kita tampilkan error dan hentikan program
                loads, err := collector.GetLoad()
                if err != nil {
                    log.Fatal(err)
                }
 
                // membuat line protocol dengan menyertakan data loads
                datapoint := fmt.Sprintf("loadavg,hostname=server-1 load1=%.2f,load5=%.2f,load15=%.2f", loads["load1"], loads["load5"], loads["load15"])
                // hit api untuk memasukan data loads kedalam database gomondb 
                data, err := command.ExecCurlInsert("-i", "http://localhost:8086/write?db={your_db_influx}&u={your_username}&p={your_password}", "--data-binary", datapoint)
                if(err == nil){
                    log.Println(data)
                }else{
                    log.Fatal(err)
                    break
                }
                // beri selang waktu 15 detik sebelum melakukan hal yang sama
                time.Sleep(15 * time.Second)
        }
}