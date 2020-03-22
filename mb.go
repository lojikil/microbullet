package main

import (
    "fmt"
    //"flag"
    "os"
    "time"
    "path"
    //"bufio"
    "regexp"
    "strconv"
    "strings"
)

func baseExists() bool {
    ok, err := os.UserHomeDir()

    if err != nil {
        return false
    }

    pathName := path.Join(ok, ".mubu")
    if _, err := os.Stat(pathName); os.IsNotExist(err) {
        return false
    }
    return true
}

func curExists() bool {
    ok, err := os.Getwd()

    if err != nil {
        return false
    }

    pathName := path.Join(ok, ".mubu")
    if _, err := os.Stat(pathName); os.IsNotExist(err) {
        return false
    }
    return true
}

func makeBase() bool {
    ok, err := os.UserHomeDir()

    if err != nil {
        return false
    }

    pathName := path.Join(ok, ".mubu")
    if _, err := os.Stat(pathName); os.IsNotExist(err) {
        if err := os.Mkdir(pathName, 0700); err != nil {
            return false
        }
        return true
    }
    return true
}

func todayExists(useCur bool) bool {

    pathName, err := getBasePath(useCur)

    if err != nil {
        return false
    }

    if _, err := os.Stat(pathName); os.IsNotExist(err) {
        return false
    }

    return true
}

func makeToday(useCur bool) bool {

    pathName, err := getBasePath(useCur)

    if err != nil {
        return false
    }

    if _, err := os.Stat(pathName); os.IsNotExist(err) {
        // this is technically a vulnerability, because
        // MkDirAll doesn't actually check each member
        // of the path for existence, and it just returns
        // nil if the path exists. I'm _roughly_ ok with
        // this for our purposes here, but we may want to
        // change this going forward
        if err := os.MkdirAll(pathName, 0700); err != nil {
            return false
        }
        return true
    }

    return true
}

func getBasePath(useCur bool) (string, error) {
    var ok string

    if useCur {
        lok, err := os.Getwd()
        if err != nil {
            return lok, err
        }

        ok = lok
    } else {
        lok, err := os.UserHomeDir()

        if err != nil {
            return lok, err
        }

        ok = lok
    }

    t := time.Now()

    // would be nice to make this work for other cases,
    // like when we want to check if a specific day in
    // the future exists
    year := fmt.Sprintf("%d", t.Year())
    month := fmt.Sprintf("%02d", t.Month())
    day := fmt.Sprintf("%02d", t.Day())
    pathName := path.Join(ok, ".mubu", year, month, day)

    fmt.Printf("In getBasePath, pathName is: %s\n", pathName)

    return pathName, nil
}

func getLatest(path string) (int, error) {
    fd, err := os.Open(path)
    if err != nil {
        return -1, err
    }

    defer fd.Close()

    files, err := fd.Readdir(-1)
    if err != nil {
        return -1, err
    }

    // sort the names, and return the latest one
    res := -1
    for _, file := range(files) {
        fmt.Printf("in loop, file: %s\n", file.Name())
        if matched, err := regexp.MatchString("^[0-9]+$", file.Name()); err == nil {
            fmt.Printf("here matched? %v\n", matched)
            // the *one* time that strconv.Atoi isn't a problem
            // because I actually want to compare on arch-non-specific
            // integers
            cur, err := strconv.Atoi(file.Name())
            if err == nil && cur > res {
                res = cur
            }
        } else {
            fmt.Printf("regexp returned error? %v", err)
        }

    }
    fmt.Printf("In getLatest, res is: %d\n", res)
    return res + 1, nil
}

func addNote(args []string, useCur bool) bool {
    var msg string

    pathName, err := getBasePath(useCur)
    if err != nil {
        return false
    }

    latestID, err := getLatest(pathName)
    if err != nil {
        return false
    }

    latest := path.Join(pathName, fmt.Sprintf("%d", latestID))

    // if we have args, use those as the note
    // otherwise, read the args from the console
    if len(args) > 0 {
        msg = fmt.Sprintf("- %s\n", strings.Join(args, " "))
    } else {
        buf := make([]byte, 8192)
        fmt.Printf("- ")
        _, err := os.Stdin.Read(buf)
        if err != nil {
            return false
        }
        msg = string(buf)
    }
    fmt.Printf("here? 192\n")

    fd, err := os.Create(latest)
    if err != nil {
        fmt.Printf("error: %v\n", err)
        return false
    }

    fmt.Printf("Here? 198\n")
    ret, err := fd.WriteString(msg)
    if err != nil || ret != len(msg) {
        return false
    }

    fd.Sync()
    fd.Close()

    return true
}

func main() {

    var useCur bool
    /*
     * I could see the case for collapsing these, and just
     * allowing all methods to operate off the user's specified
     * directory, or some sort of defaults...
     * Currently, we operate:
     * - if there is a `.mubu` in `.`, use that
     * - else, check if there is a `$HOME/.mubu`
     * - make a base
     */
    if curExists() {
        fmt.Println("using local notes repository")
        // operate out of the current directory first...
        useCur = true
    } else if baseExists() {
        // operate out of base if none exist...
        useCur = false
    } else {
        fmt.Println("no microbullet repository found")
        useCur = false
    }

    if len(os.Args) < 2 {

        if todayExists(useCur) {
            fmt.Println("today does exist...")
        } else {
            fmt.Println("No tasks for notes for today...")
        }

        os.Exit(0)
    }

    switch os.Args[1] {
        case "note", "n":
            if makeToday(useCur) {
                addNote(os.Args[2:], useCur)
            } else {
                fmt.Println("an error occurred adding today's repo...")
                os.Exit(1)
            }
        case "entry", "e":
            fmt.Println("adding a full entry...")
        case "task", "t":
            fmt.Println("adding a task...")
        case "todo", "d":
            fmt.Println("adding a todo...")
        case "init", "i":
            /*
             * I'm definitely thinking some thoughts here...
             * for one, I'd like it if we could organize multiple
             * user's notes, like say in a github repo for an assessment.
             * $PATH/.mubu/sedwards/...
             * so, thoughts:
             * - `-u` for single user, in $HOME/.mubu by default
             * - `-m` for multi-user, in the specified path  (defaults to ".")
             */
            fmt.Println("initializing the repo")
        case "header", "h":
            fmt.Println("adding a header")
        case "help", "H", "?":
            fmt.Println("printing some help...")
        case "code", "c":
            fmt.Println("adding some code...")
        case "view", "v":
            fmt.Println("viewing some notes...")
        default:
            fmt.Println("invalid command...")
    }
}
