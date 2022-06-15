### Steps to configure Talisman

Talisman is a tool that helps to scan through files to check if any sensitive data lies in your files and alerts you while commiting.

1. curl --silent  https://raw.githubusercontent.com/thoughtworks/talisman/master/global_install_scripts/install.bash > /tmp/install_talisman.bash && /bin/bash /tmp/install_talisman.bash
2. .pre-commit-config.yaml
3. pre-commit install -f
4. .talismanrc