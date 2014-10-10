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

## Contributing

1. Fork it ( https://github.com/mizoR/coo/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request
