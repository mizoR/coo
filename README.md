# coo

Manage the transcripts of ssh

## Installation and Usage

```sh
$ go get github.com/mizoR/coo/cmd/{coo,coo-tee}

$ vm_stat 3 | coo-tee vmstat.log -t

$ coo www.example.com
# Runs
#   $ ssh www.example.com | coo-tee "$HOME/.ssh/transcripts/www.example.com/YYYY-mm-dd.txt" -t -a
```
