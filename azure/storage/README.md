# Liftbox

Liftbox watches a folder for changes and react to these changes.

## Building

    $ cd azure/storage
    $ go build

## Usage

When the `config.toml` file is in the same directory of the executable `./storage`:

    $ ./storage

When the `config.toml` is in a different location:

    $ ./storage --cfg=/home/username/liftbox/config_prod.toml

When the configuration is taken from environment variables:

    $ export LIFTBOX_ROOTPATH=/home/username/liftbox/
    $ ./storage