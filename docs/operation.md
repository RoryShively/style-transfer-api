# Operation


## Make commands

```
# Start project
make start

# Run tests, select any of the below (all_tests runs e2e and unit)
make unit_tests
make e2e_tests
make all_tests

# Stop project (Deletes minikube and unmounts direectories)
make clean
```

## Curl commands

Start from project directory

```
# Start project
make start

# Get env var for API URL to run through commands
source ./hack/test-vars.sh
```

Look at available models to train on. These models represent different artists styles. The `make start` command will download the Van Gogh neural net to start with, more are available at https://hcicloud.iwr.uni-heidelberg.de/index.php/s/XXVKT5grAquXNqi

Download new models if you want, move to the `models` dir and unzip with `tar -zxvf model.tar.gz`

```
curl $API_URL/models -v

# Select model from list an set an env var for it
STYLE="van-gogh"
```

Upload an image to run through style transfer neural net. Can take about 1 minute for image to process. Make note of the id sent back to query status of image and receive stylized image later
```
curl -v -F style=$STYLE -F image=@./data/test.jpg \
        $API_URL/image

# Get id from returned json and set as env var
ID=<from returned json>
```

Check on image to see what it's status is
```
curl $API_URL/image/$ID -v
```

One the staus returns as done the stylized image can be retrieved and opened
```
curl $API_URL/image/stylized/$ID > ~/stylized.jpg && open ~/stylized.jpg
```




