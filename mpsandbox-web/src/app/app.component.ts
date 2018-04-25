import { Component } from '@angular/core';

import { WebsocketService } from './services/websocket.service';
import { InhabitantsService, Message } from './services/inhabitants.service';

import { Point, Guid, Projectile, Player, Shield } from './models/Point';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
  providers: [ WebsocketService, InhabitantsService ]
})
export class AppComponent {

  players: Player[] = [];
  player: Player = <Player> ({
      Point: new Point(50, 50),
      GUID: Guid.newGuid(),
      Name: "",
      Health: 100,
  });
  projectiles: Projectile[] = [];
  shields: Shield[] = [];
  keys: { [ keycode: string ]: boolean } = {};
  tickTimer: any;
  shieldCharge: number = 100;

  constructor(private inhabitants: InhabitantsService) {

    inhabitants.messages.subscribe(msg => {
      switch(msg.type) {
        case 'playersUpdate':
          var temp = <Player[]>JSON.parse(msg.value);
          var temp2 = [];
          for (let player of temp) {
            if (this.player.GUID !== player.GUID) {
              temp2.push(player);
            } else {
              this.player.Health = player.Health;
            }
          }
          this.players = temp2;
          break;
        case 'projectilesUpdate':
          this.projectiles = <Projectile[]>JSON.parse(msg.value);
          break;
        case 'shieldsUpdate':
          this.shields = <Shield[]>JSON.parse(msg.value);
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
          if (this.player.Health > 0) {
            switch(key) {
              case 'w':
                this.player.Point.y -= 10;
                break;
              case 'a':
                this.player.Point.x -= 10;
                break;
              case 's':
                this.player.Point.y += 10;
                break;
              case 'd':
                this.player.Point.x += 10;
                break;
              default:
                break;
            }
          }
        }
      }
      if (changed) {
        this.playerUpdate();
      }
    }, 50);

    this.playerUpdate();

  }

  playerUpdate() {
    const message = <Message>({
      type: 'playerUpdate',
      value: JSON.stringify(this.player),
    });
    this.inhabitants.messages.next(message);
  }
  
  keydown($event) {
    this.keys[$event.key] = true;
  }
  
  keyup($event) {
    this.keys[$event.key] = false;
  }
  
  mouseup($event) {
    console.log($event);
    var mouseEvent = <MouseEvent>$event;
    switch (mouseEvent.button) {
      case 0:
        if (this.player.Health > 0) {
          var fromPoint: Point = <Point>{};
          fromPoint.x = this.player.Point.x;
          fromPoint.y = this.player.Point.y;
          var toPoint: Point = <Point>{};
          toPoint.x = mouseEvent.clientX;
          toPoint.y = mouseEvent.clientY;
          var playerShootMessage = { "fromPoint": fromPoint, "toPoint": toPoint };

          const message = <Message>({
            type: 'playerShoot',
            value: JSON.stringify(playerShootMessage),
          });
          this.inhabitants.messages.next(message);
        }
      break;
      case 2:
        $event.preventDefault();
        if (this.player.Health > 0 && this.shieldCharge === 100) {
          this.shieldCharge = 0;
          this.chargeShield();
          var mouseEvent = <MouseEvent>$event;
          var point: Point = <Point>{};
          point.x = this.player.Point.x - 10;
          point.y = this.player.Point.y + 30;

          const message = <Message>({
            type: 'playerShield',
            value: JSON.stringify(point),
          });
          this.inhabitants.messages.next(message);
        }
      break;
    }
  }

  chargeShield() {
    setTimeout(()=> {
      debugger;
      this.shieldCharge += 10;
      if (this.shieldCharge < 100) {
        this.chargeShield();
      }
    }, 100);
  }
  
  contextmenu($event) {
    $event.preventDefault();
  }
  
  blur($event) {
    this.keys = {};
  }

}
