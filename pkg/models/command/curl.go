package command

import (
    "os/exec"
)

func ExecCurlInsert(optioncommand string, uri string, data string) (string, error){
		//declare command curl
	err := exec.Command("curl", optioncommand, uri, "--data-binary", data).Run()
	if err != nil {
	    return "gagal exec", err
	}
	 	
	//return data string & error nil jika command curl berhasil di running
	return "success input : "+data, nil
}

