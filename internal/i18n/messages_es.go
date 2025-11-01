// Package i18n
// Description: Spanish translations
// (c) 2024 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package i18n

func getSpanishMessages() map[MessageKey]string {
	return map[MessageKey]string{
		// Common messages
		MsgErrorLoadingConfig:   "Error al cargar la configuración: %v",
		MsgErrorReadingFile:     "Error al leer el archivo: %v",
		MsgErrorParsingYAML:     "Error al analizar YAML: %v",
		MsgErrorCreatingSymlink: "Error al crear enlace simbólico de %s a %s: %v",
		MsgErrorRemovingSymlink: "Error al eliminar enlace simbólico en %s: %v",
		MsgErrorCheckingFile:    "Error al verificar archivo en %s: %v",
		MsgErrorReadingSymlink:  "Error al leer enlace simbólico en %s: %v",
		MsgFileNotFound:         "Archivo no encontrado: %s",
		MsgInvalidBoolValue:     "Valor booleano inválido. Use 'true' o 'false'",
		MsgInvalidOS:            "SO inválido '%s'. Las opciones válidas son: linux, darwin, windows",

		// Symlinks messages
		MsgSymlinkFileNotFound:    "Archivo de enlaces simbólicos no encontrado: %s\nPor favor cree el archivo o actualice la configuración con: sok config symlinkfile <ruta>",
		MsgReadingSymlinksFrom:    "Leyendo configuración de enlaces simbólicos desde: %s",
		MsgFoundConfigurations:    "Se encontraron %d configuración(es) de enlaces simbólicos",
		MsgDryRunWouldCreate:      "[SIMULACIÓN] Se crearía enlace simbólico: %s -> %s",
		MsgSymlinkAlreadyExists:   "El enlace simbólico ya existe y es correcto: %s -> %s",
		MsgExistingSymlinkRemoved: "Enlace simbólico existente eliminado: %s",
		MsgSymlinkCreated:         "Enlace simbólico creado: %s -> %s",
		MsgSymlinkNotFound:        "Enlace simbólico no encontrado (ya eliminado): %s",
		MsgNotSymlink:             "Advertencia: %s no es un enlace simbólico, omitiendo (no se eliminarán archivos regulares)",
		MsgSymlinkWrongTarget:     "El enlace simbólico apunta a un origen diferente: %s -> %s (esperado: %s), omitiendo",
		MsgDryRunWouldRemove:      "[SIMULACIÓN] Se eliminaría enlace simbólico: %s -> %s",
		MsgSymlinkRemoved:         "Enlace simbólico eliminado: %s -> %s",
		MsgUninstallSummary:       "Resumen de Desinstalación",
		MsgWouldRemove:            "Se eliminarían: %d enlace(s) simbólico(s)",
		MsgRemoved:                "Eliminados: %d enlace(s) simbólico(s)",
		MsgNotFound:               "No encontrados: %d enlace(s) simbólico(s)",
		MsgNotSymlinks:            "No son enlaces simbólicos (omitidos): %d archivo(s)",
		MsgSkipped:                "Omitidos: %d enlace(s) simbólico(s)",
		MsgSymlinksStatus:         "Estado de Enlaces Simbólicos:",
		MsgStatus:                 "Estado",
		MsgTarget:                 "Destino",
		MsgSource:                 "Origen",
		MsgSummary:                "Resumen:",
		MsgInstalledCorrectly:     "Instalados correctamente:    %d",
		MsgWrongTarget:            "Destino incorrecto:          %d",
		MsgNotInstalled:           "No instalados:               %d",
		MsgRegularFileExists:      "Archivo regular existe:      %d",
		MsgTotalSymlinks:          "Total de enlaces simbólicos configurados: %d",
		MsgLegend:                 "Leyenda:",
		MsgLegendInstalled:        "✅ = Enlace simbólico instalado y apunta al origen correcto",
		MsgLegendWrongTarget:      "⚠️  = Enlace simbólico existe pero apunta a un origen diferente",
		MsgLegendNotInstalled:     "❌ = Enlace simbólico no instalado",
		MsgLegendRegularFile:      "⛔ = Archivo regular existe en la ubicación de destino (no es un enlace simbólico)",

		// Apply messages
		MsgApplyingChanges:    "Aplicando cambios de configuración...",
		MsgConfigReloaded:     "Configuración recargada desde el disco",
		MsgChangesToApply:     "Cambios a Aplicar",
		MsgToCreate:           "Para Crear (%d):",
		MsgToUpdate:           "Para Actualizar (%d):",
		MsgAlreadyCorrect:     "Ya Correctos: %d",
		MsgNoChangesNeeded:    "No se necesitan cambios - ¡todos los enlaces simbólicos están actualizados!",
		MsgDryRunNoChanges:    "[SIMULACIÓN] No se realizaron cambios",
		MsgApplyingChangesNow: "Aplicando Cambios",
		MsgCreated:            "Creado: %s -> %s",
		MsgUpdated:            "Actualizado: %s -> %s",
		MsgFailed:             "Fallidos: %d enlace(s) simbólico(s)",
		MsgApplySummary:       "Resumen de Aplicación",
		MsgConfigApplied:      "¡Configuración aplicada exitosamente!",

		// Init messages
		MsgInitializing:        "Inicializando sokru...",
		MsgCreatedDirectory:    "Directorio creado: %s",
		MsgDirectoryExists:     "El directorio ya existe: %s",
		MsgCreatedConfigFile:   "Archivo de configuración creado: %s",
		MsgConfigFileExists:    "El archivo de configuración ya existe: %s",
		MsgCreatedSymlinksFile: "Archivo de enlaces simbólicos creado: %s",
		MsgSymlinksFileExists:  "El archivo de enlaces simbólicos ya existe: %s",
		MsgInitSummary:         "Resumen de Inicialización",
		MsgNextSteps:           "Próximos Pasos",
		MsgAddDotfiles:         "1. Agregue sus dotfiles a los directorios de shell:",
		MsgEditSymlinks:        "2. Edite la configuración de enlaces simbólicos:",
		MsgInstallSymlinks:     "3. Instale sus enlaces simbólicos:",
		MsgViewConfig:          "4. Vea su configuración:",
		MsgInitialized:         "¡Sokru inicializado exitosamente!",

		// Config messages
		MsgConfigFile:          "Archivo de configuración: %s",
		MsgCurrentConfig:       "Configuración actual:",
		MsgDotfilesDirectory:   "Directorio de Dotfiles: %s",
		MsgSymlinksFile:        "Archivo de Enlaces:     %s",
		MsgOperatingSystem:     "Sistema Operativo:      %s",
		MsgVerbose:             "Detallado:              %v",
		MsgDryRun:              "Simulación:             %v",
		MsgCurrentDotfilesDir:  "Directorio de dotfiles actual: %s",
		MsgDotfilesDirSet:      "Directorio de dotfiles establecido en: %s",
		MsgCurrentSymlinksFile: "Archivo de enlaces simbólicos actual: %s",
		MsgSymlinksFileSet:     "Archivo de enlaces simbólicos establecido en: %s",
		MsgCurrentVerbose:      "Configuración detallada actual: %v",
		MsgVerboseSet:          "Modo detallado establecido en: %v",
		MsgCurrentDryRun:       "Configuración de simulación actual: %v",
		MsgDryRunSet:           "Modo de simulación establecido en: %v",
		MsgCurrentOS:           "Configuración de SO actual: %s",
		MsgOSSet:               "SO establecido en: %s",

		// Help messages
		MsgHelpSymlinks: "Ayuda del comando SymLinks::",
		MsgHelpConfig:   "Ayuda del comando Config::",
		MsgHelp:         "Ayuda de Sokru::",

		// Version messages
		MsgVersion: "Sok v0.1 -- HEAD",
	}
}
