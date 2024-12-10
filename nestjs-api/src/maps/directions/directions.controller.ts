import { Controller, Get, Query } from '@nestjs/common';
import { DirectionsService } from './directions.service';

@Controller('directions')
export class DirectionsController {
  constructor(private readonly directionService: DirectionsService) {}

  @Get()
  getDirections(
    @Query('originId') originId: string,
    @Query('destinationId') destinationId: string,
  ) {
    return this.directionService.getDirections(originId, destinationId);
  }
}
