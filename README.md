# logmonitor

Reads log files under /var/log/ directory and return the last events first.

Essentially it returns the original log file in reverse line order.

## Usage

Run `logmonitor` excutable file under the root directory on the target machine, and it will start a REST API server.

Endpoint: `localhost:8080/api/v1/logs`

Method: GET

Query params:

| Field  | Description | Default Value |
| ------------- | ------------- | ---- |
| filename | Log file name under the directory to query log lines for  | var5MB.txt |
| size  | Number of entries to return  | 100 |
| keyword | Filter results for log lines with keyword only | (empty, no filter) |

## Assumptions

- Each log line ends with a line break byte (`\n`) including the last line of the file.
- The maximum length of the log lines is smaller than 32KB.

## How Does It Work?

1. We read a fixed size bytes buffer from the end of the file.
2. In this buffer we locate all the `\n` bytes.
3. Between two `\n` bytes we have a line of log. Return all the log lines in reverse order.
4. Also return the location of the _first_ `\n` byte, which will serve as the starting point of the next read.
5. Repeat until we have enough lines to return.

<img width="497" alt="Screenshot 2023-08-20 at 6 04 45 PM" src="https://github.com/suyangduan/logmonitor/assets/17387788/6970b2d0-230e-428f-ba8a-9c3f5a153e14">



## For Code Reviews

[Individual commits](https://github.com/suyangduan/logmonitor/commits/main) are probably more interesting to read. Specifically

[This one](https://github.com/suyangduan/logmonitor/commit/fd84617c26874e26668a045e1c0d4fadc789c195) implements the very first operation of reading one buffer from the end of the file. (sorry i lumped in some cleaning up code 🤦)

[The follow up](https://github.com/suyangduan/logmonitor/commit/8c88ef0d5312883040336ef3b4543bb179c2e181) implements a generic version of reading buffer that takes a fileOffset so that it can read the file in any specified location.

[Then](https://github.com/suyangduan/logmonitor/commit/60d4d6e94db6df931de6547702bc56f212b8a25b) we chain up the calls to read buffer so that we can return the specified number of log lines.

The remainder is more cleaning up work. 

## Performance

There's no noticeable difference of querying 100 lines of log between log files of size 5k, 5M, and 1G (from implementation perspective it shouldn't have any difference but i should probably bench mark this just to be sure). 

For the 1GB log file which has 20M lines of logs, it takes about 12 seconds to scan through the whole thing on a macbook air.
