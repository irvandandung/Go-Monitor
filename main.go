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
    configvar, err := GetConfig("./config.yml")
    if err != nil {
        log.Fatal(err)
    }
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
        data, err := command.ExecCurlInsert("-i", UriConnection(configvar, "write"), datapoint)
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

func UriConnection(configvar *Config, paramQuery string) (uri string){
    conf := configvar.Database.Influxdb
    baseUrl := fmt.Sprintf("http://%s:%s/", conf.Host, conf.Port)
    db := fmt.Sprintf("%s", conf.Database)
    username := fmt.Sprintf("%s", conf.Username)
    password := fmt.Sprintf("%s", conf.Password)

    uri = baseUrl+paramQuery+"?db="+db+"&u="+username+"&p="+password
    return uri
}