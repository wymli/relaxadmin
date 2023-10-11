
check:
  #!/usr/bin/env bash

  go version | grep "1.20"
  if [[ "$?" != "0" ]];then echo "go 1.20 is needed";exit 1; fi

dev: check
  ENV=dev DB=0 RELOAD=1 go run main.go

build $os="linux": check
  #!/usr/bin/env bash
  set -euo pipefail

  bin="admin"
  echo "building target=$os $bin"
  rm $bin || true
  GOOS=$os go build -o $bin

download: check
  #!/usr/bin/env bash
  set -euo pipefail

  pushd static

  deps=( "https://registry.npmjs.org/vanilla-jsoneditor/-/vanilla-jsoneditor-0.18.3.tgz" 
  "https://registry.npmjs.org/jquery-ui/-/jquery-ui-1.13.2.tgz"
  "https://registry.npmjs.org/jquery/-/jquery-3.7.1.tgz"
  "https://registry.npmjs.org/jquery.json-viewer/-/jquery.json-viewer-1.5.0.tgz"
  )

  for url in ${deps[@]};
  do
    f=$(basename $url)
    f=${f%\.tgz}
    if [[ -d $f ]]; then 
      echo "skip downloading $f..."
      continue
    fi

    wget $url -O tmp.tgz
    tar -zxvf tmp.tgz
    mv package $f
    rm tmp.tgz
  done


create $resource="page" $name="":
  #!/usr/bin/env bash
  if [[ "$resource" == "page" ]] ;then
      name=${name:?"page name should be specified; just create page hello-world"}
      PAGE=$name bash scripts/just/create_page.sh > ui/page/$name.html
  fi
