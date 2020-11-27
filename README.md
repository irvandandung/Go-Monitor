# Go-Monitor
This is Golang Project for monitoring server with Curl hit &amp; database Influx

## How To Run
- First you must install influxd v1.2 & golang.
  * InfluxDb v1.2 docs : [https://archive.docs.influxdata.com/influxdb/v1.2/](https://archive.docs.influxdata.com/influxdb/v1.2/)
  * Golang docs : [https://golang.org/doc/](https://golang.org/doc/)
- Create db gomondb in influxdb.
- Clone this project.
- Go to cloned project directory.

	example :
    ```bash
       cd $your_dir/gokriyatest
    ```
- Rename config.yml.example file to config.yml
- Input the influxdb configuration in the config.yml file

	example:
    ```yaml
       database: 
         influxdb:
           host: "localhot"
           port: "8080"
           username: "root"
           password: "my_pasword12"
           database: "gomondb"
    ```
- Run Program
	
    ```bash
       go run .
    ```

- after program running, program auto insert data to database every 15 second. You can check data in you db gomondb influxdb.
  * List table auto create in db gomondb : loadavg, diskusage, memory, cpu

