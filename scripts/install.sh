#!/bin/bash
## Colors
# Black        0;30     Dark Gray     1;30
# Red          0;31     Light Red     1;31
# Green        0;32     Light Green   1;32
# Brown/Orange 0;33     Yellow        1;33
# Blue         0;34     Light Blue    1;34
# Purple       0;35     Light Purple  1;35
# Cyan         0;36     Light Cyan    1;36
# Light Gray   0;37     White         1;37

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color
 
######################################################
DIR_INSTALL_CATIGNORE=/opt/catignore/bin
LINK_INSTALL_CATIGNORE=/usr/bin

VERSION="v1.2.0"

HELP="Script to download catignore software in its latest version!

Usage:

  ${GREEN}sudo bash ${PURPLE}install.sh ${CYAN}[flags]${NC}

Flags:

  ${CYAN}--help, -h ${NC}         
    Help for catignore

   ${CYAN}--version, -v ${NC}       
    Print just the version number.

   ${CYAN}--install, -i [Options] ${NC} 
    Install from catignore on distributions: Linux

       ${CYAN}--system, -s [required]${NC} 
        Inform the supported OS
          
           ${CYAN}--architecture, -a [required]${NC} 
            Enter the architecture supported by the OS

            Example:
              ${GREEN}sudo bash ${PURPLE}install.sh ${CYAN}--install -s ${NC}linux ${CYAN}-a ${NC}amd64

   ${CYAN}--download, -d [Options]${NC} 
    Download catignore for the given operating system.

        The file is saved in the directory where the script is running.

       ${CYAN}--system, -s [required]${NC} 
      Inform the supported OS
        
         ${CYAN}--architecture, -a [required]${NC} 
          Enter the architecture supported by the OS

          Example:
            ${GREEN}bash ${PURPLE}install.sh ${CYAN}--download -s ${NC}windows ${CYAN}-a ${NC}amd64

   ${CYAN}--list, -l${NC} 
    List all OS compatible with catignore

    Example:
      ${GREEN}bash ${PURPLE}install.sh ${CYAN}--list ${NC} 

   ${CYAN}--uninstall  ${NC}        
    Uninstall catignore from linux 

    Example:
      ${GREEN}sudo bash ${PURPLE}install.sh ${CYAN}--uninstall ${NC} "


# argsCli - Get the information passed by parameter
# $1 = flags, $2 = flag -s, $3 = SO, $4 = flag -a and $5 = architecture
# Example: $1 = --install, $2 = -s, $3 = linux, $4 = -a and $5 = amd64
function argsCli() {

  [ -z "$1" ] && echo -e "${CYAN}Offer any option, use -h for help${NC}" && exit 1

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
      echo -e " ${RED} Erro: unknown ${1} flag! Type -h or --help to check the flags.${NC}"
      exit 2
  esac
}

# argsInstallDownloaded - Get the information passed by parameter, to install catignore
# $1 = flag -i, $2 = flag -s, $3 = SO, $4 = flag -a and $5 = architecture
# Example: $1 = --install, $2 = -s, $3 = linux, $4 = -a and $5 = amd64
function argsInstallDownloaded() {

  if [ -z "$2" ] || [ -z "$3" ];
  then  
    echo -e " ${YELLOW} Warning: Flag -s and OS name is required, use -h for help${NC}" 
    exit 1
  fi
  
  if [ -z "$4" ] || [ -z "$5" ];
  then  
    echo -e " ${YELLOW} Warning: Flag -a and architecture name is required, use -h for help${NC}" 
    exit 1
  fi

  argsCliArch $4 $5

  case $2 in   
    -s | --system) # Inform the OS
      checkSO $1 $3 $5
      exit 0
    ;;   
    *) # Erro unknown flags
      echo -e " ${RED} Erro: unknown ${1} flag! Use the -s flag and OS name to download.${NC}"
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
        echo -e " ${YELLOW} Warning: Architecture $2 not found for download${NC}"  
        exit 3
      fi          
    ;;
    *) # Erro unknown flags
      echo -e " ${RED} Erro: unknown ${1} flag! Use the -a flag and architecture name to download.${NC}"
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
      echo -e " ${YELLOW} Warning: There is no architecture version ${3} for ${2}${NC}"
      exit 2 
    fi
   
    if [[ "$1" == "--install" ]] ||  [[ "$1" == "-i" ]];
    then     
      echo -e " ${YELLOW} Warning: We still don't have catignore installation on OS ${2}.${NC}"
    elif [[ "$1" == "--download" ]] ||  [[ "$1" == "-d" ]];
    then
      donwloadCatIgnore $2 $3 
    fi

    exit 0
  ;;
  "darwin")

    if [[ "$3" == "386" ]]; then
      echo -e " ${YELLOW} Warning: There is no architecture version ${3} for darwin${NC}"
      exit 2 
    fi
   
    if [[ $1 == "--install" ]] ||  [[ $1 == "-i" ]];
    then     
      echo -e " ${YELLOW} Warning: We still don't have catignore installation on OS ${2}.${NC}"
    elif [[ $1 == "--download" ]] ||  [[ $1 == "-d" ]];
    then
      donwloadCatIgnore $2 $3 
    fi

    exit 0
  ;;
  *) # Erro unknown flags
    echo -e " ${RED} Erro: unknown ${1} flag! Type -h or --help to check the flags.${NC}"
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

  echo -e "${GREEN} Info: Installing catignore, please wait....${NC}"

  tar xzf catignore-*
  rm -rf catignore-*

  ln -sf $DIR_INSTALL_CATIGNORE/catignore $LINK_INSTALL_CATIGNORE/catignore

  echo -e "${GREEN} Info: Installation successfull!${NC}"

  catignore -v

  exit 0 
}

#requireSudo - Check if the script is running as sudo
function requireSudo() {
  if [[ $UID != 0 ]]; then
    echo -e " ${YELLOW} Warning: Please run this script with sudo:${NC}"
    echo -e " ${GREEN}sudo bash ${PURPLE}$0 $* ${NC}"
    exit 1
  fi
}

# donwloadCatIgnore - Download the program in the latest version
# $1 = SO and $2 = architecture
# Example: $1 = linux and $2 = amd64
function donwloadCatIgnore() {
  local SO=$1
  local ARCHITECTURE=$2

  echo -e " ${GREEN}Info: Downloading the latest version of the program for the system ${SO}-${ARCHITECTURE}\n${NC}"
 
  local DOWNLOAD=$(wget -qO- https://api.github.com/repos/miquelis/catignore/releases/latest \
  | grep browser_download_url \
  | grep $SO-$ARCHITECTURE \
  | cut -d '"' -f 4)

  wget -t 3 $DOWNLOAD
 
  if [  $? -ne 0 ]; then
    echo -e " ${RED} Erro: Failed to download the latest version of the program for the system ${SO}-${ARCHITECTURE}\n${NC}"

    echo -e " ${YELLOW} Warning: Removing installation, please wait...\n${NC}"
    uninstallLinux
    exit 1
  fi 
  
  local DOWNLOAD_NAME=$(echo $DOWNLOAD | cut -d"/" -f 9)
  
  echo -e "${GREEN} Info: Download successfull!${NC}"
  echo -e "${GREEN} Info: $PWD/${DOWNLOAD_NAME} has been saved successfully!${NC}"
}

# uninstallLinux - Removes catignore
# Don't need arguments
function uninstallLinux(){

  requireSudo

  echo -e "${GREEN} Info: Uninstall catignore, please wait...${NC}"

  rm -rf $DIR_INSTALL_CATIGNORE/catignore

  unlink $LINK_INSTALL_CATIGNORE/catignore
  
  if [ $? -ne 0 ]; then
    echo -e " ${RED} Erro: Unistall catignore failed, try again later!${NC}"
    exit 1
  fi 

  echo -e "${GREEN} Info: Uninstall successfull! Close your terminal and open it again.${NC}"

  exit 0
}

function listVersions() {

  local LIST_DOWNLOAD
  
  LIST_DOWNLOAD+=($(wget -qO- https://api.github.com/repos/miquelis/catignore/releases/latest \
  | grep browser_download_url \
  | cut -d '"' -f 4))

  
  if [ -z $LIST_DOWNLOAD ]; then
    echo -e " ${RED} Erro: Failed to list catignore versions, try again later!${NC}"
    exit 1
  fi 

  echo -e "\n ${CYAN} List of latest versions:${NC}\n"

  for value in "${LIST_DOWNLOAD[@]}"
  do
      local version=$(echo $value | cut -d"/" -f 9 | cut -d"." -f4,5 --complement)
      echo -e "-> $version"
  done 

  exit 0
}
# Start of script execution
argsCli $1 $2 $3 $4 $5