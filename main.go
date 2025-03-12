package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"
	"runtime"
)

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

	fmt.Printf("Now restart your PC/VM to have the changes take effect\n")
}

func execDocker(volume *string) {
	displayEnv := os.Getenv("DISPLAY")
	cngPathEnv := os.Getenv("CNG_PATH")

	var mountedPath *string
	if len(*volume) > 0 {
		mountedPath = volume
	} else {
		mountedPath = &cngPathEnv
	}

	user, err := user.Current()
	var userUid string
	var userGid string
	if err != nil {
		fmt.Printf("[ERROR] Could not get current user, using default user instead. Error: %a", err)
		userUid = "1000"
		userGid = "1000"
	} else {
		userUid = user.Uid
		userGid = user.Gid
	}

	cmd := exec.Command("docker",
		"run", "--name", "contiker", "-it", "--rm",
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
		"contiker/contiki-ng", "/bin/bash")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	errRun := cmd.Run()
	if errRun != nil {
		fmt.Printf("[ERROR] Command did end in error: %a\n", errRun)
	}
}

func execSh() {
	user, err := user.Current()
	var userUid string
	var userGid string
	if err != nil {
		fmt.Printf("[ERROR] Could not get current user, using default user instead. Error: %a", err)
		userUid = "1000"
		userGid = "1000"
	} else {
		userUid = user.Uid
		userGid = user.Gid
	}

	cmd := exec.Command("docker",
		"exec", "--user", fmt.Sprintf("%s:%s", userUid, userGid), "-it", "contiker", "/bin/bash")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	errRun := cmd.Run()
	if errRun != nil {
		fmt.Printf("[ERROR] Command did end in error: %a\n", errRun)
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
	shSet := flag.NewFlagSet("sh", flag.ExitOnError)
	rmSet := flag.NewFlagSet("rm", flag.ExitOnError)

	initSet := flag.NewFlagSet("init", flag.ExitOnError)
	initUrlPtr := initSet.String("git", "", "Git clone url")
	initVolumePtr := initSet.String("v", "contiki-ng", "Place to put Contiki folder")

	fixSet := flag.NewFlagSet("fix", flag.ExitOnError)
	xhostPtr := fixSet.Bool("xhost", false, "Fix xhost (X11 connectivity) issue")
	dockerPermPtr := fixSet.Bool("docker", false, "Fix Docker permission issue")

	volumePtr := flag.String("v", "", "Volume to be mounted")

	for _, v := range os.Args[1:] {
		if v == "-h" {
			fmt.Printf("Valid subcommands are: sh, rm, init, fix\n")

			flag.Usage()
			shSet.Usage()
			rmSet.Usage()
			initSet.Usage()
			fixSet.Usage()

			return
		}
	}

	// If no arguments
	if len(os.Args) == 1 {
		flag.Parse()
		execDocker(volumePtr)
	}

	// If at least one argument
	switch os.Args[1] {
	case "sh":
		shSet.Parse(os.Args[2:])
		execSh()
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
		fmt.Printf("All fixes applied\n")
	default:
		flag.Parse()
		execDocker(volumePtr)
	}
}
