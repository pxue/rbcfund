RBCFund CLI tool
================

This is a tool to easily pull and analyze openly available fund data on the RBC
investment website. RBC's websites and tools are hardly usable for beginners,
and even my advisors have a difficult time navigating their sites/tools. This
CLI tool allows me to quickly pull up freely available monthly data on any fund,
construct portfolios, display historical data, and quickly compare fund options.

# Requirements

    - Go 1.7+
    - Python2.7

# Installation

    $ go get -u github.com/pxue/rbcfund

# Usage

    $ rbcfund [<flags>] <command> [<args> ...]

# TODOs:

    * Fully configurable. ie. start and end month/year, reinvestment options,
      balance combinations
    * Better data presentation.
    * Portfolio diversity diagnosis (over/underweights)
    * Fund analysis based on actual asset content
    * Realtime monitoring

# License

Licensed under the [MIT License](./LICENSE).
