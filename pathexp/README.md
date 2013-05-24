pathexp
=======

pathexp is a library to match url path

## Installation

`go get github.com/sdming/pathexp `

## Usage

`import "github.com/sdming/pathexp" `

## Example

    pattern := "/query/hoho/{type}/{year}-{month}-{day}/"
    input := "/query/hoho/type/year-month-day/"

    re, err := pathexp.Compile(pattern)
    if err != nil {
        fmt.Println("pattern compile fail: %s; pattern=%s \n", err, pattern)
        return
    }

    args := re.FindAllStringSubmatch(input)
    for i, arg := range args {
        fmt.Println(i, arg[0], "=", arg[1])
    }
 

For more example usage, please see `pathexp_test.go`

## profiler

    (pprof) top10         
    Total: 206 samples
          82  39.8%  39.8%      187  90.8% github.com/sdming/pathexp.(*Pathex).execute
          19   9.2%  49.0%       92  44.7% runtime.mallocgc
          14   6.8%  55.8%       14   6.8% scanblock
          12   5.8%  61.7%       12   5.8% runtime.markallocated
          10   4.9%  66.5%       43  20.9% runtime.makeslice
          10   4.9%  71.4%       20   9.7% sweep
           8   3.9%  75.2%        8   3.9% MCentral_Alloc
           7   3.4%  78.6%       23  11.2% runtime.MCache_Alloc
           7   3.4%  82.0%       10   4.9% runtime.MCache_Free
           5   2.4%  84.5%        5   2.4% runtime.memmove

maybe can run more faster

## License

MIT