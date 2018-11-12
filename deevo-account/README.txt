- Build library
  - go build -buildmode=plugin
- Config in file core
  - chaincode:
    systemPlugins:
      - enabled: true
        name: deevo-account
        path: /opt/lib/deevo-account.so
        invokableExternal: true
        invokableCC2CC: true
    system:
      deevo-account: enable