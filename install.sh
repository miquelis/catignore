#!bin/bash

VERSION="v0.0.0"

HELP="Script to download catignore software in its latest version!

Usage:
  sudo bash install.sh [flags]

Flags:
  -h, --help          help for catignore
  -v, --version       Print just the version number.
  -s, --system        Tell the operating system to download. The systems are linux, windows and darwin. When informing the Linux system, it will be automatically installed on your system and if it is other systems, the download will be carried out in the directory where the script was downloaded.
  It is mandatory to use the -a or –architecture flag together with the -s, --system, otherwise an error will be generated.
  -a, --architecture   Report the architecture of the operating system. Supported architectures: linux [386, amd64, arm64], windows [386, amd64] and darwin [amd64, arm64].

Usage example install linux: sudo bash install.sh -s linux -a amd64
Usage example download: bash install.sh -s windows -a amd64"


# argsCli - Get the information passed by parameter
# $1 = flags, $2 = SO, $3 = flag -a and $4 = architecture
# Example: $1 = -s, $2 = linux, $3 = -a and $4 = amd64
function argsCli() {

  [ -z "$1" ] && echo "Offer any option, use -h for help" && exit 1

  case $1 in
    -h | --help) ## Show program help
      echo "${HELP}"
      exit 0
    ;;
    -v | --version) ## display the program version
      echo "version ${VERSION}"
      exit 0
    ;;
    -s | --system) ## Inform the operating system
      checkSO $2 $3 $4
      exit 0
    ;;
    *) ## Erro unknown flags
      echo "Erro: unknown ${1} flag! Type -h or --help to check the flags."
      exit 2
  esac
}

# checkSO - Check if the operating system is correct
# $1 = SO, $2 = flag -a and $3 = architecture
# Example: $1 = linux, $2 = -a and $3 = amd64
function checkSO() {
  
  if [ -z "$1" ] || [ -z "$2" ] || [ -z "$3" ];
  then  
    echo "Warning: Flags or fields not informed, use -h for help" 
    exit 1
  fi

  case $1 in
  "linux")

    local retval=$( argsCliArch $2 $3)

    if [[ $retval == "true" ]];then
        installCatIgnore $1 $3
    fi
    
    exit 0
  ;;
  "windows")  
   
    if [[ "$3" == "arm64" ]]; then
      echo "Warning: There is no architecture version ${3} for Windows"
      exit 2 
    fi

    local retval=$( argsCliArch $2 $3)

    if [[ $retval == "true" ]]; then
      donwloadCatIgnore $1 $3
      exit 0
    else 
      echo "$retval"
      exit 2
    fi
  ;;
  "darwin")

    if [[ "$3" == "386" ]]; then
      echo "Warning: There is no architecture version ${3} for darwin"
      exit 2 
    fi

    local retval=$( argsCliArch $2 $3)

    if [[ $retval == "true" ]];then
        donwloadCatIgnore $1 $3
    fi
   
    exit 0
  ;;
  *) ## Erro unknown flags
    echo "Erro: unknown ${1} flag! Type -h or --help to check the flags."
    exit 2
  esac

}

# argsCliArch - Get the information passed by parameter
# $1 = flag -a and $2 = architecture
# Example: $1 = -a and $2 = amd64
function argsCliArch() {

  case $1 in
    -a | --architecture)

      local retval=$( checkArchiteture $2)

      if [[ $retval == "true" ]]; then
        echo "true"
      else 
        echo "Warning: Architecture $2 not found for download"  
        exit 3
      fi          
    ;;
    *) ## Erro unknown flags
      echo "Erro: unknown ${1} flag! Type -h or --help to check the flags."
      exit 3
  esac
}

# checkArchiteture - Check if the architeture is correct
# $1 = architecture
# Example: $1 = amd64
function checkArchiteture() {
 case $1 in
  "386")
    echo "true"
  ;;
  "amd64")
    echo "true"
  ;;
  "arm64")
    echo "true"
  ;;
  *)
    echo "false"
  esac
}

#installCatIgnore - After downloading the logs, the program will be installed.
# $1 = SO and $2 = architecture
# Example: $1 = linux and $2 = amd64
function installCatIgnore() {

  requireSudo

  local PATHC=/opt/catignore
  
  mkdir -p $PATHC 
  cd $PATHC

  donwloadCatIgnoreLinux $1 $2

  if ! grep -xq "catignore.sh" /etc/profile.d/*.sh &> /dev/null; then
    local SAVEFILE="/etc/profile.d/catignore.sh"
    echo "alias catignore=$PATHC/catignore" > $SAVEFILE
    source $SAVEFILE
  fi

  echo "Info: Installation successful! Close your terminal and open it again."

  exit 0
 
}

#requireSudo - Check if the script is running as sudo
function requireSudo() {
  if [[ $UID != 0 ]]; then
    echo "Warning: Please run this script with sudo:"
    echo "sudo bash $0 $*"
    exit 1
  fi
}

#donwloadCatIgnoreLinux - Download the program in the latest version
# $1 = SO and $2 = architecture
# Example: $1 = linux and $2 = amd64
function donwloadCatIgnoreLinux() {
  local SO=$1
  local ARCHITECTURE=$2

  echo -e "Info: Downloading the latest version of the program for the system ${SO}-${ARCHITECTURE} \n"
 
  curl -s https://api.github.com/repos/miquelis/catignore/releases/latest \
  | grep browser_download_url \
  | grep $SO-$ARCHITECTURE | cut -d '"' -f 4 | xargs wget -O -  | tar -xz

  local retval="$?"

  if [ $retval -ne 0 ]; then
    echo "Warning: Failed to download the latest version of the program for the system ${SO}"
    exit 1
  fi 
}

#donwloadCatIgnore - Download the program in the latest version
# $1 = SO and $2 = architecture
# Example: $1 = linux and $2 = amd64
function donwloadCatIgnore() {
  local SO=$1
  local ARCHITECTURE=$2

  echo -e "Info: Downloading the latest version of the program for the system ${SO}-${ARCHITECTURE}" 

  local DOWNLOAD_URL=$(curl -s https://api.github.com/repos/miquelis/catignore/releases/latest | grep browser_download_url | grep $SO-$ARCHITECTURE | cut -d '"' -f 4)
  curl -LJO $DOWNLOAD_URL

  local retval="$?"

  if [ $retval -ne 0 ]; then
    echo "Warning: Failed to download the latest version of the program for the system ${SO}-${ARCHITECTURE}"
    exit 1
  fi
}


# Start of script execution
argsCli $1 $2 $3 $4

# source /etc/bashrc  &> /dev/null
# source ~/.bashrc  &> /dev/null
# source ~/.bash_profile  &> /dev/null