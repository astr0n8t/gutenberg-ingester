# gutenberg-ingester

Mirror all of [Project Gutenberg](https://www.gutenberg.org/) and keep it up to date.

Has the ability to filter based on file format and language.

See the config.yaml file for an example configuration

> [!WARNING]
> This project is provided as is and your mileage may vary.  Please submit any issues that you encounter.

## docker-compose

```yaml
gutenberg-ingester:
  image: ghcr.io/astr0n8t/gutenberg-ingester:latest
  container_name: gutenberg_ingester
  restart: always
  environment:
# optional
    - GUTENBERG_INGESTER_GUTENBERG_MIRROR_URL="https://www.gutenberg.org/" 
  volumes:
# where to store the state file
    - gutenberg_db:/var/gutenberg-ingester
# where to place the books (they can be post-processed and moved out of here)
    - /mnt/books:/data
# optional
    - /etc/gutenberg-ingester/config.yaml:/etc/gutenberg-ingester/config.yaml
```
