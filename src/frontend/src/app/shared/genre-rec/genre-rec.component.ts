import { Component } from '@angular/core';

@Component({
  selector: 'app-genre-rec',
  standalone: true,
  imports: [],
  templateUrl: './genre-rec.component.html',
  styleUrl: './genre-rec.component.css'
})
export class GenreRecComponent {
  genres: string[] = ['Action', 'Animation', 'Comedy', 'Drama', 'Horror', 'Musical', 'Children', 'All'];


}
