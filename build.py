#!/usr/bin/env python3

import os
import sys
import logging as log
import subprocess
# from threading import Thread
# import asyncio as aio

log.basicConfig(format="[%(levelname)s] %(asctime)s %(name)s: %(message)s",
                    datefmt="%Y-%m-%d %H:%M:%S",
                    level=log.INFO)

if not sys.version_info >= (3, 8):
    log.error("this script requires python 3.8 or above")
    sys.exit(1)

BIN  = "dohc"
GO   = "go"

ARCH = [ "386", "amd64", "arm", "arm64" ]
OS   = [ "linux", "android", "windows", "darwin" ]

STDOUT = sys.stdout
STDERR = sys.stderr


def system(cmd: str):
    return subprocess.run(cmd, shell=True, stdout=STDOUT, stderr=STDERR)


def build(arch: str, osname: str) -> None:
    outName = f"{BIN}-{osname}-{arch}"
    envVars = f"GOOS={osname} GOARCH={arch}"
    log.info(f"building {outName}")

    if osname == "windows":
        envVars += " CGO_ENABLED=1"
        outName += ".exe"

    
    system(f"{envVars} {GO} build -ldflags='-buildid=' -o 'build/{outName}' .")


def main() -> None:
    # loop = aio.new_event_loop()
    # aio.set_event_loop(loop)
    # thread  = Thread(target=loop.run_forever, daemon=False)
    # thread.start()

    os.makedirs("build")
    for osName in OS:
        for arch in ARCH:
            if osName == "android" and arch != "arm64":
                continue
            if osName == "darwin" and (arch == "386" or arch == "arm"):
                continue
            build(arch, osName)
            # aio.run_coroutine_threadsafe(build(arch, osName), loop)
    # thread.join()
    log.info(f"Done!")


if __name__ == "__main__":
    main()
