How to communicate to Docker Host from Docker container
=======================================================
## Solution

Map `docker.sock` as a volume from the host to the container:

## Example

- Map the volume
```bash
dock run -it --rm -v /var/run/docker.sock:/var/run/docker.sock <DKR_IMAGE> /bin/bash
```

- Test it
```bash
echo -e "GET /images/json HTTP/1.0\r\n" | nc -U /var/run/docker.sock
```

Note. If `nc -U` fails, ensure you have `netcat-openbsd` installed and not `netcat-traditional`


How to create devices from Vagrants
============
## Solution
[Link](https://gist.github.com/leifg/4713995)
