import { Component, Input, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { take } from 'rxjs';
import {
  AddSMSProviderTwilioRequest,
  AddSMTPConfigRequest,
  UpdateSMSProviderTwilioRequest,
  UpdateSMTPConfigPasswordRequest,
  UpdateSMTPConfigPasswordResponse,
  UpdateSMTPConfigRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { DebugNotificationProvider, SMSProvider, SMSProviderConfigState } from 'src/app/proto/generated/zitadel/settings_pb';
import { AdminService } from 'src/app/services/admin.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';

import { InfoSectionType } from '../../info-section/info-section.component';
import { PolicyComponentServiceType } from '../policy-component-types.enum';
import { DialogAddSMSProviderComponent } from './dialog-add-sms-provider/dialog-add-sms-provider.component';
import { PasswordDialogComponent } from './password-dialog/password-dialog.component';

@Component({
  selector: 'cnsl-notification-settings',
  templateUrl: './notification-settings.component.html',
  styleUrls: ['./notification-settings.component.scss'],
})
export class NotificationSettingsComponent implements OnInit {
  @Input() public serviceType!: PolicyComponentServiceType;
  public smsProviders: SMSProvider.AsObject[] = [];
  public logNotificationProvider!: DebugNotificationProvider.AsObject;
  public fileNotificationProvider!: DebugNotificationProvider.AsObject;

  public smtpLoading: boolean = false;
  public smsProvidersLoading: boolean = false;
  public logProviderLoading: boolean = false;
  public fileProviderLoading: boolean = false;

  public form!: FormGroup;

  public SMSProviderConfigState: any = SMSProviderConfigState;
  public InfoSectionType: any = InfoSectionType;

  public hasSMTPConfig: boolean = false;

  // show available providers

  constructor(
    private service: AdminService,
    private dialog: MatDialog,
    private toast: ToastService,
    private fb: FormBuilder,
    private authService: GrpcAuthService,
  ) {
    this.form = this.fb.group({
      senderAddress: [{ disabled: true, value: '' }, [Validators.required]],
      senderName: [{ disabled: true, value: '' }, [Validators.required]],
      tls: [{ disabled: true, value: true }, [Validators.required]],
      host: [{ disabled: true, value: '' }, [Validators.required]],
      user: [{ disabled: true, value: '' }, [Validators.required]],
    });
  }

  ngOnInit(): void {
    this.fetchData();
    this.authService
      .isAllowed(['iam.write'])
      .pipe(take(1))
      .subscribe((allowed) => {
        if (allowed) {
          this.form.enable();
        }
      });
  }

  private fetchData(): void {
    this.smtpLoading = true;
    this.service
      .getSMTPConfig()
      .then((smtpConfig) => {
        this.smtpLoading = false;
        if (smtpConfig.smtpConfig) {
          this.hasSMTPConfig = true;
          this.form.patchValue(smtpConfig.smtpConfig);
        }
      })
      .catch((error) => {
        this.smtpLoading = false;
        if (error && error.code === 5) {
          console.log(error);
          this.hasSMTPConfig = false;
        }
      });

    this.smsProvidersLoading = true;
    this.service
      .listSMSProviders()
      .then((smsProviders) => {
        this.smsProvidersLoading = false;
        if (smsProviders.resultList) {
          this.smsProviders = smsProviders.resultList;
        }
      })
      .catch((error) => {
        this.smsProvidersLoading = false;
        this.toast.showError(error);
      });

    this.logProviderLoading = true;
    this.service
      .getLogNotificationProvider()
      .then((logNotificationProvider) => {
        this.logProviderLoading = false;
        if (logNotificationProvider.provider) {
          this.logNotificationProvider = logNotificationProvider.provider;
        }
      })
      .catch(() => {
        this.logProviderLoading = false;
      });

    this.fileProviderLoading = true;
    this.service
      .getFileSystemNotificationProvider()
      .then((fileNotificationProvider) => {
        this.fileProviderLoading = false;
        if (fileNotificationProvider.provider) {
          this.fileNotificationProvider = fileNotificationProvider.provider;
        }
      })
      .catch(() => {
        this.fileProviderLoading = false;
      });
  }

  private updateData(): Promise<UpdateSMTPConfigPasswordResponse.AsObject> | any {
    if (this.hasSMTPConfig) {
      const req = new UpdateSMTPConfigRequest();
      req.setHost(this.host?.value ?? '');
      req.setSenderAddress(this.senderAddress?.value ?? '');
      req.setSenderName(this.senderName?.value ?? '');
      req.setTls(this.tls?.value ?? false);
      req.setUser(this.user?.value ?? '');

      return this.service.updateSMTPConfig(req).catch((error) => {
        this.toast.showError(error);
      });
    } else {
      const req = new AddSMTPConfigRequest();
      req.setHost(this.host?.value ?? '');
      req.setSenderAddress(this.senderAddress?.value ?? '');
      req.setSenderName(this.senderName?.value ?? '');
      req.setTls(this.tls?.value ?? false);
      req.setUser(this.user?.value ?? '');

      return this.service.addSMTPConfig(req).catch((error) => {
        this.toast.showError(error);
      });
    }
  }

  public savePolicy(): void {
    const prom = this.updateData();
    if (prom) {
      prom
        .then(() => {
          this.toast.showInfo('SETTING.SMTP.SAVED', true);
          setTimeout(() => {
            this.fetchData();
          }, 2000);
        })
        .catch((error: unknown) => {
          this.toast.showError(error);
        });
    }
  }

  public editSMSProvider(): void {
    const dialogRef = this.dialog.open(DialogAddSMSProviderComponent, {
      width: '400px',
      data: {
        smsProviders: this.smsProviders,
      },
    });

    dialogRef.afterClosed().subscribe((req: AddSMSProviderTwilioRequest | UpdateSMSProviderTwilioRequest) => {
      if (req) {
        if (this.hasTwilio) {
          this.service
            .updateSMSProviderTwilio(req as UpdateSMSProviderTwilioRequest)
            .then(() => {
              this.toast.showInfo('SETTING.SMS.TWILIO.ADDED', true);
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        } else {
          this.service
            .addSMSProviderTwilio(req as AddSMSProviderTwilioRequest)
            .then(() => {
              this.toast.showInfo('SETTING.SMS.TWILIO.ADDED', true);
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        }
      }
    });
  }

  public setSMTPPassword(): void {
    const dialogRef = this.dialog.open(PasswordDialogComponent, {
      width: '400px',
      data: {
        i18nTitle: 'SETTING.SMTP.SETPASSWORD',
        i18nLabel: 'SETTING.SMTP.PASSWORD',
      },
    });

    dialogRef.afterClosed().subscribe((password: string) => {
      if (password) {
        const passwordReq = new UpdateSMTPConfigPasswordRequest();
        passwordReq.setPassword(password);

        this.service
          .updateSMTPConfigPassword(passwordReq)
          .then(() => {
            this.toast.showInfo('SETTING.SMTP.PASSWORDSET', true);
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
  }

  public get twilio(): SMSProvider.AsObject | undefined {
    return this.smsProviders.find((p) => p.twilio);
  }

  public get senderAddress(): AbstractControl | null {
    return this.form.get('senderAddress');
  }

  public get senderName(): AbstractControl | null {
    return this.form.get('senderName');
  }

  public get tls(): AbstractControl | null {
    return this.form.get('tls');
  }

  public get user(): AbstractControl | null {
    return this.form.get('user');
  }

  public get host(): AbstractControl | null {
    return this.form.get('host');
  }

  public get hasTwilio(): boolean {
    const twilioProvider: SMSProvider.AsObject | undefined = this.smsProviders.find((p) => p.twilio);
    if (twilioProvider && !!twilioProvider.twilio) {
      return true;
    } else {
      return false;
    }
  }
}
