#!/bin/bash

DIR_INSTALL_CATIGNORE=/opt/catignore
DIR_SERIVCE_CATIGNORE=/etc/profile.d/catignore.sh

VERSION="v1.1.3"

HELP="Script to download catignore software in its latest version!

Usage:

  sudo bash install.sh [flags]

Flags:

  --help, -h          
    Help for catignore

  --version, -v       
    Print just the version number.

  --install, -i [Options] 
    Install from catignore on distributions: Linux

      --system, -s [required]
        Inform the supported OS
          
          --architecture, -a [required]
            Enter the architecture supported by the OS

            Example:
              sudo bash install.sh --install -s linux -a amd64

  --download, -d [Options]
    Download catignore for the given operating system.

        The file is saved in the directory where the script is running.

      --system, -s [required]
      Inform the supported OS
        
        --architecture, -a [required]
          Enter the architecture supported by the OS

          Example:
            bash install.sh --download -s windows -a amd64

  --list, -l
    List all OS compatible with catignore

    Example:
      bash install.sh --list

  --uninstall         
    Uninstall catignore from linux 

    Example:
      sudo bash install.sh --uninstall"


# argsCli - Get the information passed by parameter
# $1 = flags, $2 = flag -s, $3 = SO, $4 = flag -a and $5 = architecture
# Example: $1 = --install, $2 = -s, $3 = linux, $4 = -a and $5 = amd64
function argsCli() {

  [ -z "$1" ] && echo -e "Offer any option, use -h for help" && exit 1

  case $1 in
    --help | -h) # Show program help
      echo -e "${HELP}"
      exit 0
    ;;
    --version | -v) # display the program version
      echo -e "version ${VERSION}"
      exit 0
    ;;
    --list | -l) # display the program version
      listVersions
      exit 0
    ;;
    --install | -i) # Install
      
      argsInstallDownloaded $1 $2 $3 $4 $5
      exit 0
    ;;
    --download | -d) # Download 
      
      argsInstallDownloaded $1 $2 $3 $4 $5
      exit 0
    ;;
    --uninstall) # Uninstaller catignore
      uninstallLinux
      exit 0
    ;;
    *) ## Erro unknown flags
      echo -e "Erro: unknown ${1} flag! Type -h or --help to check the flags."
      exit 2
  esac
}

# argsInstallDownloaded - Get the information passed by parameter, to install catignore
# $1 = flag -i, $2 = flag -s, $3 = SO, $4 = flag -a and $5 = architecture
# Example: $1 = --install, $2 = -s, $3 = linux, $4 = -a and $5 = amd64
function argsInstallDownloaded() {

  if [ -z "$2" ] || [ -z "$3" ];
  then  
    echo -e "Warning: Flag -s and OS name is required, use -h for help" 
    exit 1
  fi
  
  if [ -z "$4" ] || [ -z "$5" ];
  then  
    echo -e "Warning: Flag -a and architecture name is required, use -h for help" 
    exit 1
  fi

  argsCliArch $4 $5

  case $2 in   
    -s | --system) # Inform the OS
      checkSO $1 $3 $5
      exit 0
    ;;   
    *) # Erro unknown flags
      echo -e "Erro: unknown ${1} flag! Use the -s flag and OS name to download."
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

      if [[ $retval != "true" ]]; then
        echo -e "Warning: Architecture $2 not found for download"  
        exit 3
      fi          
    ;;
    *) # Erro unknown flags
      echo -e "Erro: unknown ${1} flag! Use the -a flag and architecture name to download."
      exit 3
  esac
}

# checkSO - Check if the OS is correct
# $1 = flag -i, $2 = SO and $3 = architecture
# Example: $1 = --install, $2 = linux and $3 = amd64
function checkSO() {  
  
  case $2 in
  "linux")

      if [[ "$1" == "--install" ]] ||  [[ $1 == "-i" ]];
      then     
        installCatIgnoreLinux $2 $3
      elif [[ "$1" == "--download" ]] ||  [[ $1 == "-d" ]];
      then
        donwloadCatIgnore $2 $3 
      fi

    exit 0
  ;;
  "windows")  
   
    if [[ "$3" == "arm64" ]]; then
      echo -e "Warning: There is no architecture version ${3} for ${2}"
      exit 2 
    fi
   
    if [[ "$1" == "--install" ]] ||  [[ "$1" == "-i" ]];
    then     
      echo -e "Warning: We still don't have catignore installation on OS ${2}."
    elif [[ "$1" == "--download" ]] ||  [[ "$1" == "-d" ]];
    then
      donwloadCatIgnore $2 $3 
    fi

    exit 0
  ;;
  "darwin")

    if [[ "$3" == "386" ]]; then
      echo -e "Warning: There is no architecture version ${3} for darwin"
      exit 2 
    fi
   
    if [[ $1 == "--install" ]] ||  [[ $1 == "-i" ]];
    then     
      echo -e "Warning: We still don't have catignore installation on OS ${2}."
    elif [[ $1 == "--download" ]] ||  [[ $1 == "-d" ]];
    then
      donwloadCatIgnore $2 $3 
    fi

    exit 0
  ;;
  *) # Erro unknown flags
    echo -e "Erro: unknown ${1} flag! Type -h or --help to check the flags."
    exit 2
  esac

}

# checkArchiteture - Check if the architeture is correct
# $1 = architecture
# Example: $1 = amd64
function checkArchiteture() {
 case $1 in
  "386")
    echo -e "true"
  ;;
  "amd64")
    echo -e "true"
  ;;
  "arm64")
    echo -e "true"
  ;;
  *)
    echo -e "false"
  esac
}

#installCatIgnoreLinux - After downloading the logs, the program will be installed.
# $1 = SO and $2 = architecture
# Example: $1 = linux and $2 = amd64
function installCatIgnoreLinux() {

  requireSudo

  mkdir -p $DIR_INSTALL_CATIGNORE
  cd $DIR_INSTALL_CATIGNORE

  donwloadCatIgnore $1 $2

  echo -e "Info: Installing catignore, please wait...."

  tar xzf catignore-*
  rm -rf catignore-*

  if ! grep -xq "catignore.sh" /etc/profile.d/*.sh &> /dev/null; then
    echo -e "alias catignore=$DIR_INSTALL_CATIGNORE/catignore" > $DIR_SERIVCE_CATIGNORE  
    source $DIR_SERIVCE_CATIGNORE
  fi

  echo -e "Info: Installation successfull! Close your terminal and open it again."

  exit 0 
}

#requireSudo - Check if the script is running as sudo
function requireSudo() {
  if [[ $UID != 0 ]]; then
    echo -e "Warning: Please run this script with sudo:"
    echo -e "sudo bash $0 $*"
    exit 1
  fi
}

# donwloadCatIgnore - Download the program in the latest version
# $1 = SO and $2 = architecture
# Example: $1 = linux and $2 = amd64
function donwloadCatIgnore() {
  local SO=$1
  local ARCHITECTURE=$2

  echo "Info: Downloading the latest version of the program for the system ${SO}-${ARCHITECTURE}\n"
 
  local DOWNLOAD=$(wget -qO- https://api.github.com/repos/miquelis/catignore/releases/latest \
  | grep browser_download_url \
  | grep $SO-$ARCHITECTURE \
  | cut -d '"' -f 4)

  wget -t 3 $DOWNLOAD
 
  if [  $? -ne 0 ]; then
    echo -e "Erro: Failed to download the latest version of the program for the system ${SO}-${ARCHITECTURE}\n"

    echo -e "Warning: Removing installation, please wait...\n"
    uninstallLinux
    exit 1
  fi 
  
  local DOWNLOAD_NAME=$(echo $DOWNLOAD | cut -d"/" -f 9)
  
  echo -e "Info: Download successfull!"
  echo -e "Info: $PWD/${DOWNLOAD_NAME} has been saved successfully!"
}

# uninstallLinux - Removes catignore
# Don't need arguments
function uninstallLinux(){

  requireSudo

  echo -e "Info: Uninstall catignore, please wait..."

  rm -rf $DIR_INSTALL_CATIGNORE
  rm -rf $DIR_SERIVCE_CATIGNORE
  
  if [ $? -ne 0 ]; then
    echo -e "Erro: Unistall catignore failed, try again later!"
    exit 1
  fi 

  echo -e "Info: Uninstall successfull! Close your terminal and open it again."

  exit 0
}

function listVersions() {

  local LIST_DOWNLOAD
  
  LIST_DOWNLOAD+=($(wget -qO- https://api.github.com/repos/miquelis/catignore/releases/latest \
  | grep browser_download_url \
  | cut -d '"' -f 4))

  
  if [ -z $LIST_DOWNLOAD ]; then
    echo -e "Erro: Failed to list catignore versions, try again later!"
    exit 1
  fi 

  echo -e "\nList of latest versions:"

  for value in "${LIST_DOWNLOAD[@]}"
  do
      local version=$(echo $value | cut -d"/" -f 9 | cut -d"." -f4,5 --complement)
      echo -e " -> $version"
  done 

  exit 0
}
# Start of script execution
argsCli $1 $2 $3 $4 $5