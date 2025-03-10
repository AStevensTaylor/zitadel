import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';

import { GrpcAuthService } from './grpc-auth.service';

declare const tinycolor: any;

export interface Color {
  name: string;
  hex: string;
  rgb: string;
  contrastColor: string;
}

@Injectable()
export class ThemeService {
  private _darkTheme: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(true);
  public isDarkTheme: Observable<boolean> = this._darkTheme.asObservable();

  private primaryColorPalette: Color[] = [];
  private warnColorPalette: Color[] = [];
  private backgroundColorPalette: Color[] = [];

  constructor(private authService: GrpcAuthService) {
    const theme = localStorage.getItem('theme');
    if (theme) {
      if (theme === 'light-theme') {
        this.setDarkTheme(false);
      } else {
        this.setDarkTheme(true);
      }
    }
  }

  setDarkTheme(isDarkTheme: boolean): void {
    this._darkTheme.next(isDarkTheme);
  }

  public updateTheme(colors: Color[], type: string, theme: string): void {
    colors.forEach((color) => {
      document.documentElement.style.setProperty(`--theme-${theme}-${type}-${color.name}`, color.hex);
      document.documentElement.style.setProperty(`--theme-${theme}-${type}-contrast-${color.name}`, color.contrastColor);
    });
  }

  public savePrimaryColor(color: string, isDark: boolean): void {
    this.primaryColorPalette = this.computeColors(color);
    this.updateTheme(this.primaryColorPalette, 'primary', isDark ? 'dark' : 'light');
  }

  public saveWarnColor(color: string, isDark: boolean): void {
    this.warnColorPalette = this.computeColors(color);
    this.updateTheme(this.warnColorPalette, 'warn', isDark ? 'dark' : 'light');
  }

  public saveBackgroundColor(color: string, isDark: boolean): void {
    this.backgroundColorPalette = this.computeColors(color);
    this.updateTheme(this.backgroundColorPalette, 'background', isDark ? 'dark' : 'light');
  }

  public saveTextColor(colorHex: string, isDark: boolean): void {
    const theme = isDark ? 'dark' : 'light';
    document.documentElement.style.setProperty(`--theme-${theme}-${'text'}`, colorHex);
    const secondaryTextHex = tinycolor(colorHex).setAlpha(0.78).toHex8String();
    document.documentElement.style.setProperty(`--theme-${theme}-${'secondary-text'}`, secondaryTextHex);
  }

  private computeColors(hex: string): Color[] {
    return [
      this.getColorObject(tinycolor(hex).lighten(52), '50'),
      this.getColorObject(tinycolor(hex).lighten(37), '100'),
      this.getColorObject(tinycolor(hex).lighten(26), '200'),
      this.getColorObject(tinycolor(hex).lighten(12), '300'),
      this.getColorObject(tinycolor(hex).lighten(6), '400'),
      this.getColorObject(tinycolor(hex), '500'),
      this.getColorObject(tinycolor(hex).darken(6), '600'),
      this.getColorObject(tinycolor(hex).darken(12), '700'),
      this.getColorObject(tinycolor(hex).darken(18), '800'),
      this.getColorObject(tinycolor(hex).darken(24), '900'),
      this.getColorObject(tinycolor(hex).lighten(50).saturate(30), 'A100'),
      this.getColorObject(tinycolor(hex).lighten(30).saturate(30), 'A200'),
      this.getColorObject(tinycolor(hex).lighten(10).saturate(15), 'A400'),
      this.getColorObject(tinycolor(hex).lighten(5).saturate(5), 'A700'),
    ];
  }

  private getColorObject(value: any, name: string): Color {
    const c = tinycolor(value);
    return {
      name: name,
      hex: c.toHexString(),
      rgb: c.toRgbString(),
      contrastColor: this.getContrast(c.toHexString()),
    };
  }

  public isLight(hex: string): boolean {
    const color = tinycolor(hex);
    return color.isLight();
  }

  public isDark(hex: string): boolean {
    const color = tinycolor(hex);
    return color.isDark();
  }

  public getContrast(color: string): string {
    const onBlack = tinycolor.readability('#000', color);
    const onWhite = tinycolor.readability('#fff', color);
    if (onBlack > onWhite) {
      return 'hsla(0, 0%, 0%, 0.87)';
    } else {
      return '#ffffff';
    }
  }

  public setDefaultColors = () => {
    const darkPrimary = '#bbbafa';
    const lightPrimary = '#5469d4';

    const darkWarn = '#ff3b5b';
    const lightWarn = '#cd3d56';

    const darkBackground = '#111827';
    const lightBackground = '#fafafa';

    const darkText = '#ffffff';
    const lightText = '#000000';

    this.savePrimaryColor(darkPrimary, true);
    this.savePrimaryColor(lightPrimary, false);

    this.saveWarnColor(darkWarn, true);
    this.saveWarnColor(lightWarn, false);

    this.saveBackgroundColor(darkBackground, true);
    this.saveBackgroundColor(lightBackground, false);

    this.saveTextColor(darkText, true);
    this.saveTextColor(lightText, false);
  };

  public loadPrivateLabelling(): void {
    this.setDefaultColors();

    const isDark = (color: string) => this.isDark(color);
    const isLight = (color: string) => this.isLight(color);

    this.authService
      .getMyLabelPolicy()
      .then((lpresp) => {
        const labelpolicy = lpresp.policy;

        const darkPrimary = labelpolicy?.primaryColorDark || '#bbbafa';
        const lightPrimary = labelpolicy?.primaryColor || '#5469d4';

        const darkWarn = labelpolicy?.warnColorDark || '#ff3b5b';
        const lightWarn = labelpolicy?.warnColor || '#cd3d56';

        let darkBackground = labelpolicy?.backgroundColorDark;
        let lightBackground = labelpolicy?.backgroundColor;

        let darkText = labelpolicy?.fontColorDark ?? '#ffffff';
        let lightText = labelpolicy?.fontColor ?? '#000000';

        this.savePrimaryColor(darkPrimary, true);
        this.savePrimaryColor(lightPrimary, false);

        this.saveWarnColor(darkWarn, true);
        this.saveWarnColor(lightWarn, false);

        if (darkBackground && !isDark(darkBackground)) {
          console.info(
            `Background (${darkBackground}) is not dark enough for a dark theme. Falling back to zitadel background`,
          );
          darkBackground = '#111827';
        }
        this.saveBackgroundColor(darkBackground || '#111827', true);

        if (lightBackground && !isLight(lightBackground)) {
          console.info(
            `Background (${lightBackground}) is not light enough for a light theme. Falling back to zitadel background`,
          );
          lightBackground = '#fafafa';
        }
        this.saveBackgroundColor(lightBackground || '#fafafa', false);

        if (darkText && !isLight(darkText)) {
          console.info(
            `Text color (${darkText}) is not light enough for a dark theme. Falling back to zitadel's text color`,
          );
          darkText = '#ffffff';
        }
        this.saveTextColor(darkText || '#ffffff', true);

        if (lightText && !isDark(lightText)) {
          console.info(
            `Text color (${lightText}) is not dark enough for a light theme. Falling back to zitadel's text color`,
          );
          lightText = '#000000';
        }
        this.saveTextColor(lightText || '#000000', false);
      })
      .catch((error) => {
        console.error('could not load private labelling policy!', error);
        this.setDefaultColors();
      });
  }
}
