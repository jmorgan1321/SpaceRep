$(document).ready(function() {
    $('.flashcard').on('click', function() {
        $('.flashcard').toggleClass('flipped');
    });

    $('#AcceptBtn').on('click', function() {
        var data = {
            Status: "Accept",
        }

        $.ajax({
            type: "POST",
            url: "/api/v1/review/accept",
            data: JSON.stringify(data),
            contentType: 'application/json; charset=utf-8',
            dataType: 'json',
        })
        .done(function(data, textStatus, xhr) {
            console.log("done");
            console.log("textStatus: " + textStatus);
            console.dir(data);
            console.dir(xhr);

            $('.flashcard').toggleClass('flipped', false);
            $('#flashcardDiv').html(data.newCard);
        });
    });

    $('#ForgotBtn').on('click', function() {
        var data = {
            Status: "Forgot",
        }

        $.ajax({
            type: "POST",
            url: "/api/v1/review/forgot",
            data: JSON.stringify(data),
            contentType: 'application/json; charset=utf-8',
            dataType: 'json',
        })
        .always(function(a, textStatus, b) {
            console.log("finished");
            console.log("textStatus: " + textStatus);
            console.dir(a);
            console.dir(b);

            $('.flashcard').toggleClass('flipped', false);
            $('#flashcardDiv').html(data.newCard);
        });
    });

    $('#SaveBtn').on('click', function() {
        var data = {
            Status: "Save",
        }

        $.ajax({
            type: "POST",
            url: "/api/v1/save",
            data: JSON.stringify(data),
            contentType: 'application/json; charset=utf-8',
            dataType: 'json',
        })
        .always(function(a, textStatus, b) {
            console.log("finished");
            console.log("textStatus: " + textStatus);
            console.dir(a);
            console.dir(b);
        });

        alert("Saved Session Data")
    });


});
