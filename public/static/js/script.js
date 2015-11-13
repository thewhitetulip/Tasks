/*

This is a javascript file for omninotesweb
============

Author:  Suraj patil
Updated: January 2015
keyCode: n-110
*/

$(document).ready(function(){

  /*on() is used instead of click because click can be used only on static elements, and on() is to be used when you add
  elements dynamically*/
  $('[data-toggle="tooltip"]').tooltip();
//    
//    $('.items').openOnHover(function(){
//        alert();
//    });

  $('#addNote').click(function(){
    $('#addNoteModal').modal('show');
  });

  $('.floating-action-icon-add').click(function(){
    $('#addNoteModal').modal('show');
  });

   //$(document).on("click", '.note-close',closeDelete); //when you delete a note, the x on the top right corner

   $(document).on('click','.open-note', openNote); //when you want to open a note in full screen, the second icon on the bottom right corner from the right

   //$(document).on('click','.hashtag', hashTag); //function to handle search by hashtag *TODO*

   //$('#addNoteModalSaveBtn').click(addNoteToDOM); //Adds note to the DOM

   /*$( document ).keypress(
     function(event){
       if ( event.which == 110 ) { //bind the 'n' key to add note
           $('#addNoteModal').modal('show');
        }

        if (event.which==109  ) { //binds the 'm' key to show the navigation drawer
          $('.sidebar-toggle').click();
        }
     }
   );*/

});

function addNoteToDOM(){
  var title = $('#add-note-title').val();
  var content = $('#add-note-content').val();
  if (title!="" || content!=""){
    var note=$('<div class="col-md-4 col-sm-4 "><div class="panel note-sm"><div class="panel-heading" >'+title +'<button class="close note-close" > &times;</button></div><hr style="margin:0 0 3px;"><div class="panel-body">'+content+' </div></div></div>');
    $('.col-md-12.row').prepend(note);
    $('#addNoteModal').modal('hide');
    $('#add-note-title').val("");
    $('#add-note-content').val("");
  }
  else{
    alert("Empty note can't be saved!");
  }
}

function closeDelete(){
  var note = $(this).parent().parent().parent();
  note.fadeOut('slow');
  note.remove();

}

function openNote(){
  var element = $(this);
  var cont=element.parent().parent().siblings().contents().toArray();
  var note_body =cont[2].data;
  console.log(note_body);
  var note_title = cont[0].data;
  console.log(note_title);


  var ONmodal = $('#openNoteModal');
  ONmodal.find('.modal-title').text(note_title);
  ONmodal.find('.modal-body').text(note_body);
  ONmodal.modal('show');
}

function hashTag(){
  alert($(this).html());
}
