package main
 
import (
    "encoding/json"
    "fmt"
    "log"
    "time"
    // import package collector sesuai dengan module path
    "github.com/irvandandung/gomon/collector"
    // import package command sesuai dengan module path
    "github.com/irvandandung/gomon/pkg/models/command"
    //import package gocron
    "github.com/go-co-op/gocron"
)
 
func main() {
    configvar, err := GetConfig("./config.yml")
    if err != nil {
        log.Fatal(err)
    }


    //inisialisasi gocron
    cron := gocron.NewScheduler(time.Local)


    //load running code every 15 second
    cron.Every(15).Second().Do(func() {
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
            cron.Clear()
        }        
    })


    //storage running code every 15 second
    cron.Every(15).Second().Do(func() {
        // panggil collector.GetStorage() yang akan mengembalikan 2 nilai
        // map[string]string dan error, map[string]string yang berisi
        // mountpoint dengan informasi mountpoint tersebut akan kita simpan
        // di variable storages
        storages, err := collector.GetStorage()
        if err != nil {
            log.Fatal(err)
        }     

        for storage := range storages {
            // buat s sebagai alias StorageInfo yang berada pada package
            // collector
            var s collector.StorageInfo
            // lakukan unmarshal pada value map, dan simpan hasilnya pada
            // StorageInfo, yang sebelumnya kita beri alias s
            if err := json.Unmarshal([]byte(storages[storage]), &s); err != nil {
                log.Fatal(err)
            }
           // buat variable datapoint yang berisi nama mountpoint dan informasi
            // mengenai total ukuran storage, kapasitas yang tersedia, dan
            // kapasitas yang sudah digunakan
            datapoint := fmt.Sprintf("diskusage,hostname=server-1,mountpoint=%s size=%d,avail=%d,used=%d", storage, s.Size, s.Avail, s.Used)
            // hit api untuk memasukan data diskusage kedalam database gomondb 
            data, err := command.ExecCurlInsert("-i", UriConnection(configvar, "write"), datapoint)
            if(err == nil){
                log.Println(data)
            }else{
                log.Fatal(err)
                cron.Clear()
            }
        }  
    })


    //memory running code every 15 second
    cron.Every(15).Second().Do(func() {
        // panggil fungsi GetMemory() yang berada pada package collector
        // fungsi tsb akan mengembalikan 2 nilai, map[string]uint64 dan error
        // map[string]uint64 yang berisi informasi mengenai memori akan
        // disimpan pada variable memory
        memory, err := collector.GetMemory()
        if err != nil {
            log.Fatal(err)
        }
 
        // buat line protocol yang berisi informasi mengenai memori
        memdp := fmt.Sprintf("memory,hostname=server-1 memtotal=%d,memused=%d", memory["MemTotal"], memory["Used"])
        // hit api untuk memasukan data memory kedalam database gomondb 
        data, err := command.ExecCurlInsert("-i", UriConnection(configvar, "write"), memdp)
        if(err == nil){
            log.Println(data)
        }else{
            log.Fatal(err)
            cron.Clear()
        }        
    })
    

    //start all cron in code
    cron.StartBlocking()
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