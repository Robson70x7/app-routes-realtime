import { SubscribeMessage, WebSocketGateway } from '@nestjs/websockets';
import { RoutesService } from '../routes.service';

@WebSocketGateway({
  cors: {
    origin: '*',
  },
})
export class RoutesDriverGateway {
  constructor(private readonly routeService: RoutesService) {}

  @SubscribeMessage('client:new-points')
  async handleMessage(client: any, payload: any) {
    const { route_id } = payload;
    const route = await this.routeService.findOne(route_id);
    // @ts-expect-error - routes has not defined
    const { steps } = route.directions.routes[0].legs[0];
    for (const step of steps) {
      const { lat, lng } = step.start_location;

      client.emit(`server:new-points/${route_id}:list`, {
        route_id,
        lat,
        lng,
      });

      client.broadcast.emit('server:new-points:list', {
        route_id,
        lat,
        lng,
      });

      await sleep(2000);

      client.emit(`server:new-points/${route_id}:list`, {
        route_id,
        lat: step.end_location.lat,
        lng: step.end_location.lng,
      });

      client.broadcast.emit('server:new-points:list', {
        route_id,
        lat: step.end_location.lat,
        lng: step.end_location.lng,
      });
      await sleep(2000);
    }
  }
}
const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));
