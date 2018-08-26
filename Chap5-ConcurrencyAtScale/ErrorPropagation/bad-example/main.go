package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
)

type MyError struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]interface{}
}

func wrapError(err error, messagef string, msgArgs ...interface{}) MyError {
	return MyError{
		// 1. we store the error which we're wrapping. We always want to be able to
		//  get back to the lowest-level error in case we need to investigate what happened.
		Inner:   err,
		Message: fmt.Sprintf(messagef, msgArgs...),
		// 2. This line of code takes note of the stack trace when the error was created.
		//   A more sophisticated error type might elide the stack-frame from wrapError.
		StackTrace: string(debug.Stack()),
		// 3. Here we create a catch-all for storing miscellaneous information.
		//   This is where we might store the concurrent ID, a hash of the stack trace,
		//   other contextual information that might help in diagnosing the error.
		Misc: make(map[string]interface{}),
	}
}

func (err MyError) Error() string {
	return err.Message
}

// lowlevel module
type LowLevelErr struct {
	error
}

func isGloaballyExec(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		// 1. Here we wrap the raw error from calling os.Stat with a customized error.
		//   In this case we are OK with the message coming out of this error, and so we won't mask it.
		return false, LowLevelErr{(wrapError(err, err.Error()))}
	}
	return info.Mode().Perm()&0100 == 0100, nil
}

// "intermediate" module
type IntermediateErr struct {
	error
}

func runJob(id string) error {
	const jobBinPath = "/bad/job/binary"
	isExecutable, err := isGloaballyExec(jobBinPath)
	if err != nil {
		// (1) Here we are passing on errors from the lowlevel module. Because of our
		//   architectural decision to consider errors passed on from other modules without
		//   wrapping them in our own type bugs, this will cause us issues later.
		// 這裡我們並沒有用我們自己定義的Intermideate error，故印出的訊息會難以讓人讀懂
		return err
	} else if isExecutable == false {
		return wrapError(nil, "job binary is not executable")
	}
	// (1)
	return exec.Command(jobBinPath, "--id="+id).Run()
}

func handleError(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logID: %v]: ", key))
	log.Printf("%#v", err)
	fmt.Printf("[%v] %v", key, message)
}

func main() {
	log.SetOutput(os.Stdout)
	// 3. Here we log out the full error in case someone needs to dig into what happend.
	log.SetFlags(log.Ltime | log.LUTC)

	err := runJob("1")
	if err != nil 
	  // 若沒有自定義的訊息，預設一個若有錯誤發生則印出該訊息，
		msg := "There was an unexpected issue; please report this as a bug."
		// 1. Here we check to see if the error is of the expected type.
		//  If it is, we know it's a well-crafted error, and we can simply pass its message on to the user.
		if _, ok := err.(IntermediateErr); ok {
			msg = err.Error()
		}
		// 2. On this line we bind th log and error message together with an ID of 1.
		//   We could easily make this increase monotonically, or use a GUID to ensure a unique ID.
		handleError(1, err, msg)
	}
}
