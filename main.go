package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"
	"runtime"
	"strings"

	_ "github.com/SamClercky/contiker/pkg/pkgmanager"
)

type Container struct {
	Command      string `json:"Command"`
	CreatedAt    string `json:"CreatedAt"`
	ID           string `json:"ID"`
	Image        string `json:"Image"`
	Labels       string `json:"Labels"`
	LocalVolumes string `json:"LocalVolumes"`
	Mounts       string `json:"Mounts"`
	Names        string `json:"Names"`
	Ports        string `json:"Ports"`
	RunningFor   string `json:"RunningFor"`
	State        string `json:"State"`
	Status       string `json:"Status"`
}

func checkContikerUp() (bool, error) {
	cmd := exec.Command("docker", "container", "ls",
		"-f", "name=contiker",
		"--format", "json")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var containers []Container

	for _, line := range lines {
		var container Container
		if len(line) == 0 {
			continue
		}

		if err := json.Unmarshal([]byte(line), &container); err != nil {
			fmt.Println("Error parsing JSON:", err)
			continue
		}
		containers = append(containers, container)
	}

	return len(containers) > 0, nil
}

func execFixXhost() {
	fmt.Printf("Starting xhost fix\n")
	cmd := exec.Command("xhost", "+local:docker")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	errRun := cmd.Run()
	if errRun != nil {
		fmt.Printf("[ERROR] Command did end in error: %a\n", errRun)
	}

	fmt.Println("> Xhost fix successful")
}

func execFixDocker() {
	fmt.Printf("Starting Docker fix\n")
	user, err := user.Current()
	if err != nil {
		fmt.Printf("[ERROR] Could not get user: %a\n", err)
		os.Exit(-1)
	}
	cmd := exec.Command("sudo", "usermod", "-aG", user.Username, "docker")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	errRun := cmd.Run()
	if errRun != nil {
		fmt.Printf("[ERROR] Command did end in error: %a\n", errRun)
	}

	fmt.Println("> Docker fix successful")
	fmt.Printf("Now restart your PC/VM to have the changes take effect\n")
}

func execFixFilePermissions() {
	fmt.Printf("Starting file permissions fix\n")
	user, err := user.Current()
	if err != nil {
		fmt.Printf("[ERROR] Could not get user: %a\n", err)
		os.Exit(-1)
	}

	isUp, err := checkContikerUp()
	if err != nil {
		fmt.Printf("[ERROR] Could not fix file permissions as it was not possible to check if Contiki was up with error: %a\n", err)
		return
	}
	if !isUp {
		fmt.Printf("[ERROR] Could not fix file permissions as there was no Contiki instance up, please run the following command:\n\n\tcontiker\n\n")
		return
	}

	cmd := exec.Command("docker", "exec", "-it", "contiker", "chown", "-R", fmt.Sprintf("%s:%s", user.Uid, user.Gid), "/home/user/contiki-ng")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	errRun := cmd.Run()
	if errRun != nil {
		fmt.Printf("[ERROR] Command did end in error: %a\n", errRun)
	}

	fmt.Println("> File permission fix successful")
}

func execDocker(volume *string, startCmd string, isRoot bool) {
	displayEnv := os.Getenv("DISPLAY")
	cngPathEnv := os.Getenv("CNG_PATH")

	var mountedPath *string
	if len(*volume) > 0 {
		mountedPath = volume
	} else {
		mountedPath = &cngPathEnv
	}

	var userUid string
	var userGid string

	if isRoot {
		userUid = "0"
		userGid = "0"
	} else {
		user, err := user.Current()
		if err != nil {
			fmt.Printf("[ERROR] Could not get current user, using default user instead. Error: %a", err)
			userUid = "1000"
			userGid = "1000"
		} else {
			userUid = user.Uid
			userGid = user.Gid
		}
	}

	// check if up
	isUp, err := checkContikerUp()
	if err != nil {
		fmt.Printf("Could not check if contiker is up with error: %a\n", err)
		os.Exit(-1)
	}

	if isUp && len(*volume) > 0 {
		fmt.Println("[WARN] You are specifying a volume while a container is already activated. The previous setup will remain. To create a new contiker environment with the specified volume, run:\n\tcontiker rm && contiker -v", *volume)
		fmt.Println()
	}

	if !isUp {
		cmd := exec.Command("docker",
			"run", "--name", "contiker", "-it", "--rm", "-d",
			"--privileged",
			"--ipc=host",
			"--network", "host",
			"-e", fmt.Sprintf("DISPLAY=%s", displayEnv),
			"-e", "_JAVA_AWT_WM_NONREPARENTING=1",
			"-e", fmt.Sprintf("LOCAL_UID=%s", userUid), "-e", fmt.Sprintf("LOCAL_GID=%s", userGid),
			"-e", "JDK_JAVA_OPTIONS='-Dawt.useSystemAAFontSettings=on -Dswing.aatext=true -Dswing.defaultlaf=com.sun.java.swing.plaf.gtk.GTKLookAndFeel -Dsun.java2d.opengl=true'",
			"-v", "/dev/:/dev/",
			"-v", "/tmp/.X11-unix:/tmp/.X11-unix",
			"--mount", fmt.Sprintf("type=bind,source=%s,destination=/home/user/contiki-ng", *mountedPath),
			"contiker/contiki-ng")

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout

		errRun := cmd.Run()
		if errRun != nil {
			fmt.Printf("[ERROR] Could not start contiker with error: %a\n", errRun)
			fmt.Printf("If the error was something related to having not enough permissions and docker, then try to run:\n\tcontiker fix -docker\n\n")
		}
	}

	cmd := exec.Command("docker",
		"exec", "--user", fmt.Sprintf("%s:%s", userUid, userGid), "-it", "contiker", startCmd)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	errRun := cmd.Run()
	if errRun != nil {
		fmt.Printf("[ERROR] Command did end in error: %a\n", errRun)
	}
	if errRun != nil && startCmd == "cooja" {
		fmt.Printf("If the error was something related to X11 and %s, then try to run:\n\tcontiker fix -xhost\n\n", displayEnv)
	}
}

func execRm() {
	cmd := exec.Command("docker",
		"container", "rm", "-f", "contiker")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	errRun := cmd.Run()
	if errRun != nil {
		fmt.Printf("[ERROR] Command did end in error: %a\n", errRun)
	}
}

func execInit(url string, folder string) {
	if len(url) == 0 {
		url = "https://github.com/contiki-ng/contiki-ng.git"
	}

	numCpus := runtime.NumCPU()
	cmd := exec.Command("git",
		"clone",
		"--recurse-submodules",
		fmt.Sprintf("-j%d", numCpus),
		url,
		folder)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	errRun := cmd.Run()
	if errRun != nil {
		fmt.Printf("[ERROR] Command did end in error: %a\n", errRun)
	}

	// Print some help information
	cwd, _ := os.Getwd()
	fmt.Printf("\n\nTo make the current directory the permanent Contiki install, add the following line to your .bashrc\n")
	fmt.Printf("export CNG_PATH=\"%s\"\n\n", path.Join(cwd, folder))
}

func main() {
	upSet := flag.NewFlagSet("up", flag.ExitOnError)
	rmSet := flag.NewFlagSet("rm", flag.ExitOnError)

	initSet := flag.NewFlagSet("init", flag.ExitOnError)
	initUrlPtr := initSet.String("git", "", "Git clone url")
	initVolumePtr := initSet.String("v", "contiki-ng", "Place to put Contiki folder")

	fixSet := flag.NewFlagSet("fix", flag.ExitOnError)
	xhostPtr := fixSet.Bool("xhost", false, "Fix xhost (X11 connectivity) issue")
	dockerPermPtr := fixSet.Bool("docker", false, "Fix Docker permission issue")
	filePermPtr := fixSet.Bool("fileperm", false, "Fix file permission issues by changing the ownership of all files in the current Contiki instance to the current user")

	volumePtr := flag.String("v", "", "Volume to be mounted")
	execPtr := flag.String("e", "/bin/bash", "Run command")
	rootPtr := flag.Bool("root", false, "Execute as root")
	resetPtr := flag.Bool("reset", false, "Stop already running Contiker instances before starting")

	for _, v := range os.Args[1:] {
		if v == "-h" {
			fmt.Println("Valid subcommands are: up, rm, init, fix, cooja")
			fmt.Println()

			fmt.Println("== General running commands ==")
			flag.Usage()
			fmt.Println()

			fmt.Println("== UP: check docker up status ==")
			upSet.Usage()
			fmt.Println()

			fmt.Println("== RM: Remove docker container ==")
			rmSet.Usage()
			fmt.Println()

			fmt.Println("== INIT: Setup Contiki environment ==")
			initSet.Usage()
			fmt.Println()

			fmt.Println("== FIX: Fix common issues ==")
			fixSet.Usage()

			return
		}
	}

	startDocker := func() {
		flag.Parse()

		if *resetPtr {
			up, err := checkContikerUp()
			if err != nil {
				fmt.Printf("[ERROR] Could not check if contiker was up, with error: %a\n", err)
			} else {
				if up {
					execRm()
				}
				// If not up, then we don't need to remove the container
			}
		}

		execDocker(volumePtr, *execPtr, *rootPtr)
	}

	// If no arguments
	if len(os.Args) == 1 {
		startDocker()
		return
	}

	// If at least one argument
	switch os.Args[1] {
	case "up":
		upSet.Parse(os.Args[2:])
		isUp, err := checkContikerUp()
		if err != nil {
			fmt.Printf("Error while checking if contiker was up, with error: %a\n", err)
			os.Exit(-1)
		}

		if isUp {
			fmt.Printf("Contiker is up. To shut it down run `contiker rm`\n")
		} else {
			fmt.Printf("Contiker is down.\n")
		}
	case "rm":
		rmSet.Parse(os.Args[2:])
		execRm()
	case "init":
		initSet.Parse(os.Args[2:])
		execInit(*initUrlPtr, *initVolumePtr)
	case "fix":
		fixSet.Parse(os.Args[2:])
		if *xhostPtr {
			execFixXhost()
		}
		if *dockerPermPtr {
			execFixDocker()
		}
		if *filePermPtr {
			execFixFilePermissions()
		}
		fmt.Printf("All fixes applied\n")
	case "cooja":
		coojaStr := "cooja"
		execPtr = &coojaStr
		startDocker()
	default:
		startDocker()
	}
}
