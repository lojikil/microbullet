# Overview

MicroBullet is a simple command line journal for note taking; it's mostly a playground for me to play with Go a bit, so the code 
is currently *atrocious*. I want to make it simple to keep a daily journal for myself that can easily be placed on the web, so
currently it uses the file system as a daily note & task storage. Eventually, I'll probably use a simple database and then "export"
to the file system.

**WARNING** the code is a mess. I'm just playing around a bit, and know I need to clean things up.

# Usage

- `note` (or `n`): add a note, either to the current directory's `.mubu` repo, or your user's homedir
