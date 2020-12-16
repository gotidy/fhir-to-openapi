package main

import "flag"

func ParseFlags() Config {
	config := Config{}
	flag.StringVar(&(config.Output), "o", "", "Output file, else output to STDOUT")
	flag.StringVar(&(config.Input), "i", "", "Input file, else get from STDIN")
	flag.Parse()
	return config
}
