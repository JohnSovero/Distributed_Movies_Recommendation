import { Component } from '@angular/core';
import { BannerComponent } from '../../shared/banner/banner.component';
import { GenreRecComponent } from '../../shared/genre-rec/genre-rec.component';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [BannerComponent, GenreRecComponent], // Include HttpClientModule here
  templateUrl: './home.component.html',
  styleUrl: './home.component.css'
})
export class HomeComponent {
  
}
