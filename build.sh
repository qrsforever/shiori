## !/bin/bash

# SHIORI_DIR=/blog/source/_shiori SOCKS5_PROXY=127.0.0.1:1080 /data/shiori/shiori serve -p 8698

top_dir=`dirname ${BASH_SOURCE[0]}`

dbroot="/blog/source/_shiori"
image="shiori"

result=$(docker images --format "{{.ID}}" $image)

if [[ x$result == x || x$1 == x1 ]]
then
    cd $top_dir
    docker build --tag $image --file Dockerfile .
    cd - > /dev/null
fi

# docker run -dit --name shiori --network bridge \
#     --restart unless-stopped --env SOCKS5_PROXY=127.0.0.1:1080 \
#     --publish 8698:8080 \
#     --volume $dbroot:/srv/shiori \
#     --entrypoint /bin/sh $image \
#     -c "shiori serve -p 8080"
