import { Controller, Get, Query } from '@nestjs/common';
import { PlacesService } from './places.service';

@Controller('places')
export class PlacesController {
  constructor(private readonly placeService: PlacesService) {}

  @Get()
  findPlaces(@Query('text') text: string) {
    return this.placeService.findPLaces(text);
  }
}
