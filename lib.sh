#!/bin/zsh

# Check if a commit message was provided
if [ -z "$1" ]; then
  echo "Error: No module version provided."
  echo "Usage: $0 <SemVer>"
  exit 1
fi

# Validate manual version if provided
if [ ! -z "$1" ]; then
  if [[ ! $1 =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: Provided version is not a valid SemVer format (vMAJOR.MINOR.PATCH)."
    exit 1
  fi
fi

url="github.com/matthxwpavin/ticketing@$1"
rule="\n🦇👻🖤💀🎃👽🔮⚰️🧛‍♂️🧟‍♂️🦉🌌🦇🧛‍♀️🌙🖤🦴👽🔫🧟‍♀️🦇🖤👻🕶️🧛‍♂️🕯️💀🦠⏳🦉🦴🦇👻🧛‍♀️🧟🖤🦴🧛‍♂️🦉🕷️🖤🧟‍♀️⚰️\n"
upgrade="echo '${rule}' && pwd && go get ${url}"
cd ./auth && eval ${upgrade}
cd ../tickets && eval ${upgrade}
cd ../orders && eval ${upgrade}
cd ../expiration && eval ${upgrade}
cd ../payment && eval ${upgrade}
printf ${rule}
