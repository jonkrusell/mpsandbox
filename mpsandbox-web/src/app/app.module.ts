import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { NgModule } from '@angular/core';

import { WebsocketService } from './services/websocket.service';
import { InhabitantsService } from './services/inhabitants.service';

import { AppComponent } from './app.component';

@NgModule({
  declarations: [
    AppComponent,
  ],
  imports: [
    BrowserModule,
    FormsModule,
  ],
  providers: [ InhabitantsService, WebsocketService ],
  bootstrap: [ AppComponent ],
})
export class AppModule { }
