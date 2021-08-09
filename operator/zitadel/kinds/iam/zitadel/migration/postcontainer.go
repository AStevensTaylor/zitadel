package migration

import (
	"github.com/caos/zitadel/operator/helpers"
	"strings"

	"github.com/caos/zitadel/operator/common"

	corev1 "k8s.io/api/core/v1"
)

func getPostContainers(
	dbHost string,
	dbPort string,
	migrationUser string,
	secretPasswordName string,
	customImageRegistry string,
	version string,
	dbCerts string,
) []corev1.Container {

	return []corev1.Container{
		{
			Name:    "delete-flyway-user",
			Image:   common.BackupImage.Reference(customImageRegistry, version),
			Command: []string{"/bin/bash", "-c", "--"},
			Args: []string{
				strings.Join([]string{
					deleteUserCommand(envMigrationUser, deleteFile),
					"cockroach.sh sql --certs-dir=" + certTempMountPath + " --host=" + dbHost + ":" + dbPort + " -e \"$(cat " + deleteFile + ")\";",
				}, ";"),
			},
			Env: baseEnvVars(envMigrationUser, envMigrationPW, migrationUser, secretPasswordName),
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      dbCerts,
					MountPath: dbCerts,
				},
			},
			SecurityContext: &corev1.SecurityContext{
				RunAsUser:    helpers.PointerInt64(1000),
				RunAsGroup:   helpers.PointerInt64(1000),
				RunAsNonRoot: helpers.PointerBool(true),
			},
			TerminationMessagePath:   corev1.TerminationMessagePathDefault,
			TerminationMessagePolicy: "File",
			ImagePullPolicy:          "IfNotPresent",
		},
	}
}

func deleteUserCommand(user, file string) string {
	return strings.Join([]string{
		"echo -n 'DROP USER IF EXISTS ' > " + file,
		"echo -n ${" + user + "} >> " + file,
		"echo -n ';' >> " + file,
	}, ";")
}
