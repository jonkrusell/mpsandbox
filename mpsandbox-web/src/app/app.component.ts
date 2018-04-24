import { Component } from '@angular/core';

import { WebsocketService } from './services/websocket.service';
import { InhabitantsService, Message } from './services/inhabitants.service';

import { Point, Guid } from './models/Point';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
  providers: [ WebsocketService, InhabitantsService ]
})
export class AppComponent {

  points: Point[] = [];
  player: Point = <Point> ({
      x: 50,
      y: 50,
      GUID: Guid.newGuid(),
      Name: ""
  });
  keys: { [ keycode: string ]: boolean } = {};
  tickTimer: any;

  constructor(private inhabitants: InhabitantsService) {
    // this.points.push(new Point(20, 50));
    // this.points.push(new Point(150, 100));
    // this.points.push(new Point(175, 75));

    inhabitants.messages.subscribe(msg => {			
      switch(msg.type) {
        case 'pointsUpdate':
          var temp = <Point[]>JSON.parse(msg.value);
          var temp2 = [];
          for (let point of temp) {
            if (this.player.GUID !== point.GUID) {
              temp2.push(point);
            }
          }
          this.points = temp2;
          break;
        default:
          console.log('unknown message type: ', msg);
          break;
      }
    });
    
    this.tickTimer = setInterval(() => {
      var changed = false;
      for (var key in this.keys) {
        if(this.keys[key] && this.keys[key] === true) {
          changed = true;
          switch(key) {
            case 'w':
              this.player.y -= 10;
              break;
            case 'a':
              this.player.x -= 10;
              break;
            case 's':
              this.player.y += 10;
              break;
            case 'd':
              this.player.x += 10;
              break;
            default:
              break;
          }
        }
      }
      if (changed) {
        const message = <Message>({
          type: 'playerUpdate',
          value: JSON.stringify(this.player),
        });
        inhabitants.messages.next(message);
      }
    }, 50);

  }
  
  keydown($event) {
    this.keys[$event.key] = true;
  }
  
  keyup($event) {
    this.keys[$event.key] = false;
  }

  private message = {
		type: 'update',
		value: 'this is a test message'
	}

  sendMsg() { 
		this.inhabitants.messages.next(this.message);
		this.message.value = '';
	}

}
