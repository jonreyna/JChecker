JChecker: Executes commands on Juniper devices at set time intervals.
===============================================

- [Intro](#intro)
- [Installation](#installation)
- [User Guide](#user-guide)
    - [Help Flag](#help-flag)
    - [Command Config File](#command-config-file)
- [License](#license)

Intro
-----

JChecker is a command line application that allows you to execute certain
commands on Juniper devices at user defined time intervals, until a deadline
is reached.

Commands currently supported are:

    - show chassis environment

Command execution results are output to a user defined file. Only
CSV output files are currently supported, but more formats will be
available in the future.

Installation
------------

1. Clone the source.
2. Use go build (or go install) to build it.
3. Put the binary in your path.


User Guide
----------

### Help flag
JChecker will keep the help flag up to date as changes are made.
If commands are added, they (and their flags) will have help as well.

    meddling_monk at TRex in ~
    $ JChecker --help
    JChecker executes requested commands at user defined time intervals, until
    a user defined deadline is reached. Deadlines, time intervals, and device
    definitions are defined in a user provided CSV file (the command CSV file).

    It's format is:
    command, ip, interval, deadline, results file, username, password

    A comment character '#' at the beginning of the line, causes the parser to
    ignore the entire line.

    Commands currently supported include:
      - show chassis environment

    Usage:
    JChecker [flags]

    Flags:
      -c, --command-config="jchecker_command_config.csv": Location of the file with CSV commands.
      -h, --help=false: help for JChecker

### Command Config File
JChecker needs a command config file so it knows where, when, what to execute. Command config files
also specify NETCONF credentials for the Juniper device.

An example command config file could look like this:

    # command, IP address, interval, deadline, output file, user name, password
    show chassis environment, 192.168.1.3, 1m, 15m, out_show_chassis_env.csv, root, abc123

The above config file will execute the equivalent of `show chassis environment` via NETCONF on
IP `192.168.1.3` every `1m` (one minute) for `15m` (15 minutes), and will output the results to the file
`out_show_chassis_env.csv`. It will use the user name `root` and password `abc123` as its
credentials.

License
-------
This software is licensed under the BSD 2-clause “Simplified” License
© 2015 JChecker contributors
