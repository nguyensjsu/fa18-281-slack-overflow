$(".add-cart").on("click",function(){
    var item_id = this.value;
    var restraunt_id = this.getAttribute("res_id");
    var item_name = this.getAttribute("item_name");
    var item_price = this.getAttribute("item_price");
    var restraunt_name = $("#restraunt_name").text();
    var res_name = restraunt_name.trim();
    console.log("res_id:"+this.getAttribute("res_id"))
    console.log("item_id:"+this.value)
    console.log("THIS WORKS")
    $.ajax({
      type: 'POST',
      url: 'localhost:30000/api/cart/itemSave',
      data: { 
       restraunt_id: restraunt_id, 
       item_id: item_id,
       item_name : item_name,
       item_price : item_price,
       restraunt_name : res_name 
      },
      dataType: "json",

      crossDomain: true,
      success: function (msg) {
         // console.log(msg[0])
         var element;
         $('#add-remove').empty()
         msg.forEach(function(value) {
            //do something with values
            console.log(value);
            element = '<li><span class="quantity">1x</span>'+
                 '<span class="product">'+value.item_name+'</span>'+
                 '<span class="money">$'+value.price+'</span>'+
                 '<div class="controls"><span aria-hidden="true">[<!-- </span> <a href="javascript:OLO.Order.edit(672716372);">edit</a><span aria-hidden="true"> |'+  
                 '--></span>'+
                    '<a href="javascript:OLO.Order.delete(672716372);" data-value='+value.item_id+'>remove</a> <span aria-hidden="true">]</span>'+
                 '</div>'+
                 '<div class="options"></div></li>';

          $("#bSubTotal").html(value.total)
          $("#add-remove").append(element);
         console.log(element);
        });
         /*var element = '<li><span class="quantity">1x</span>'+
                 '<span class="product">Shishito Peppers</span>'+
                 '<span class="money">$6.50</span>'+
                 '<div class="controls"><span aria-hidden="true">[<!-- </span> <a href="javascript:OLO.Order.edit(672716372);">edit</a><span aria-hidden="true"> |'+  
                 '--></span>'+
                    '<a href="javascript:OLO.Order.delete(672716372);">remove</a> <span aria-hidden="true">]</span>'+
                 '</div>'+
                 '<div class="options"></div></li>';*/
         // $("#BasketProducts").append(element);

      },
      error: function (request, status, error) {

         alert(error);
      }
   });
  });