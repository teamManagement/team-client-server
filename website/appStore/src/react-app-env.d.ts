/// <reference types="react-scripts" />
import { ElectronAPI } from '@electron-toolkit/preload'

declare global {
  //#region APP相关接口
  enum AppType {
    REMOTE_WEB,
    LOCAL_WEB
  }

  enum IconType {
    URL,
    ICON_FONT
  }

  interface AppInfo {
    id: string
    name: string
    inside: boolean
    type: AppType
    remoteSiteUrl: string
    url: string
    icon: string
    iconType: IconType
    desc: string
    shortDesc: string
    version: string
  }
  //#endregion

  //#region 网络代理相关接口
  interface HttpOptions {
    method?: 'GET' | 'POST' | 'DELETE' | 'PUT' | 'OPTION'
    jsonData?: any
    header?: { [key: string]: string }
  }

  interface ResponseError {
    error?: boolean
    httpCode: number
    code: string
    message: string
    raw: string
  }

  interface TcpTransferInfo<T> {
    cmdCode: number
    data: T
    errMsg: string
  }

  enum TcpTransferCmdCode {
    BLOCKING_CONNECTION,
    RESTORE_SERVER_ERR,
    RESTORE_SERVER_OK
  }

  interface ProxyApi {
    httpWebServerProxy<T>(url: string, options?: HttpOptions): Promise<T>
    httpLocalServerProxy<T>(url: string, options?: HttpOptions): Promise<T>
    registerServerMsgHandler<T>(fn: (data: TcpTransferInfo<T>) => void): Promise<void>
    removeServerMsgHandler(fn: (data: TcpTransferInfo<any>) => void): void
  }
  //#endregion

  //#region electronApi代理
  interface MenuItemOptions {
    role?:
      | 'undo'
      | 'redo'
      | 'cut'
      | 'copy'
      | 'paste'
      | 'pasteAndMatchStyle'
      | 'delete'
      | 'selectAll'
      | 'reload'
      | 'forceReload'
      | 'toggleDevTools'
      | 'resetZoom'
      | 'zoomIn'
      | 'zoomOut'
      | 'toggleSpellChecker'
      | 'togglefullscreen'
      | 'window'
      | 'minimize'
      | 'close'
      | 'help'
      | 'about'
      | 'services'
      | 'hide'
      | 'hideOthers'
      | 'unhide'
      | 'quit'
      | 'showSubstitutions'
      | 'toggleSmartQuotes'
      | 'toggleSmartDashes'
      | 'toggleTextReplacement'
      | 'startSpeaking'
      | 'stopSpeaking'
      | 'zoom'
      | 'front'
      | 'appMenu'
      | 'fileMenu'
      | 'editMenu'
      | 'viewMenu'
      | 'shareMenu'
      | 'recentDocuments'
      | 'toggleTabBar'
      | 'selectNextTab'
      | 'selectPreviousTab'
      | 'mergeAllWindows'
      | 'clearRecentDocuments'
      | 'moveTabToNewWindow'
      | 'windowMenu'
    type?: 'normal' | 'separator' | 'submenu' | 'checkbox' | 'radio'
    label?: string
    sublabel?: string
    toolTip?: string
    icon?: string
    enabled?: boolean
    acceleratorWorksWhenHidden?: boolean
    visible?: boolean
    checked?: boolean
    registerAccelerator?: boolean
    submenu?: MenuItemOptions[]
    id?: string
    click?: () => void
    before?: string[]
    after?: string[]
    beforeGroupContaining?: string[]
    afterGroupContaining?: string[]
  }

  interface ContextMenu {
    appendMenuItem(item: MenuItemOptions): Promise<void>
    registerItemClick(id: string, fn: () => void | any): void
    clearItems(): Promise<void>
    popup(): Promise<void>
  }

  class MenuList {
    public constructor(menuItems: MenuItemOptions[], menuId?: string)
    items(): MenuItem[] | undefined
    popup(): Promise<void>
    click(itemId: string): void
    id(): string
  }

   //#region 应用视图相关
   class ApplicationView {
    public constructor(id: string, url: string)
    public init(): Promise<void>
    public destroy(): Promise<void>
    public hide(): Promise<void>
    public load(): Promise<void>
    public show(data?: {
      x: number
      y: number
      width?: number
      widthOffset: number
      height?: number
      heightOffset: number
    }): Promise<void>
  }
  //#endregion

  interface Window {

    teamworkSDK: {
      // contextmenu: {
      //   build(menuItems: MenuItemOptions, menuId?: string): Menu
      //   clear(id: string): void
      //   clearAll(): void
      // }
      store: {
        set(key: string, val: any): Promise<void>
        get<T>(key: string, defaultValue?: any): Promise<T>
        has(key: string): Promise<void>
        delete(key: string): Promise<void>
        clear(): Promise<void>
      }
      exec: {
        lookPath(name: string): Promise<string>
        run(
          cmd: string,
          options?: { env?: [key: string]; cwd?: string }
        ): Promise<{ exitCode: number; stderr?: string; stdout?: string }>
      }
      proxy: {
        isEnabled443(): Promise<boolean>
        enable443(): Promise<void>
        disable443(): Promise<void>
      }
      cache: {
        set(): Promise<void>
        file: {
          store(localPath: string, expire?: number):Promise<string>,
          delete(cacheId: string):Promise<void>,
          getDownloadUrl(cacheId: string):string,
        }
      }
    }

    electron: ElectronAPI & {
      ContextMenu: {
        get(): ContextMenu
        getById(menuId: string): ContextMenu
      }
    }
    api: {
      login(username: string, password: string): Promise<void>
      contextmenu: {
        build(menuItems: MenuItemOptions[], menuId?: string): Menu
        clear(menuId: string): void
        clearAll(): void
      }
    }
    proxyApi: ProxyApi
    app: {
      getOpenedIdList(): string[]
      getApplicationViewById(id: string): ApplicationView | undefined
      openApp(
        id: string,
        url: string,
        bounds?: {
          x?: number
          y?: number
          width?: number
          widthOffset?: number
          height?: number
          heightOffset?: number
        }
      ): Promise<ApplicationView>
      closeApp(id: string): Promise<void>
      listenOpenStatusNotice(
        listenId: string,
        fn: (id: string, status: 'open' | 'close') => void
      ): void
      removeListenOpenStatusNotice(listenId: string): void
      show(
        id: string,
        bounds?: {
          x?: number
          y?: number
          width?: number
          widthOffset?: number
          height?: number
          heightOffset?: number
        }
      ): Promise<boolean>
      // showOrLoad(
      //   id: string,
      //   url: string,
      //   bounds?: {
      //     x?: number
      //     y?: number
      //     width?: number
      //     widthOffset?: number
      //     height?: number
      //     heightOffset?: number
      //   }
      // ): Promise<boolean>
      hangUp(): Promise<void>
      restore(): Promise<void>
      install(appId: string): Promise<void>
      uninstall(appId: string): Promise<void>
    }
  }
}
